package kademlia


import (
	"sync"
	"time"
	"fmt"
)

type File struct {
	Filename string
	Data []byte
	Pinned chan bool
}

type FileString struct {
	Filename string
	Data string
	Pinned bool
}

type DFS struct {
	RoutingTable 	*RoutingTable
	PurgeTimer 		int
	Kademlia			Kademlia
	mux 					sync.Mutex // ← this mutex protects the cache below
	Files 				map[string]File
}

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func NewDFS(RoutingTable 	*RoutingTable, PurgeTimer int) DFS{
	return DFS{
		RoutingTable: RoutingTable,
		PurgeTimer: PurgeTimer,
		Files: make(map[string]File),
	}
}

func (dfs *DFS) InitDFS(kademlia Kademlia) {
	dfs.Kademlia = kademlia
}

func (dfs *DFS) Cat(filename string, done chan []byte) {
  done<-dfs.Files[filename].Data
}

func (dfs *DFS) Store(filename string, data []byte, done chan string) {
	// Lock so only one goroutine at a time can write.
	dfs.mux.Lock()

	dfs.Files[filename] = File{filename, data, make(chan bool)}

	go dfs.PurgeFile(filename)

	dfs.mux.Unlock()

	done<-filename
}

func (dfs *DFS) Pin(filename string, done chan bool) {
	go dfs.StopPurge(filename)
  done<-true
}

func (dfs *DFS) Unpin(filename string, done chan bool) {
	go dfs.PurgeFile(filename)
  done<-true
}

func (dfs *DFS) PurgeFile(filename string) {
	select {
	case <-dfs.Files[filename].Pinned:
			break
		case <-time.After(time.Duration(dfs.PurgeTimer) * time.Millisecond):

			data := dfs.Files[filename].Data

			delete(dfs.Files, filename)

			target := NewContact(NewHashKademliaID(filename), "localhost:1111")
			contacts := dfs.Kademlia.FindClosestInCluster(&target)

			for _, c := range contacts {
				doneStore := make(chan string)
				go dfs.Kademlia.Network.SendStoreMessage(&dfs.Kademlia.RoutingTable.me,
					&c, filename, data, doneStore)
				<-doneStore
			}
	}
}

func (dfs *DFS) StopPurge(filename string) {
	dfs.Files[filename].Pinned <- true
}

func (dfs *DFS) GetFiles() []FileString {
	var r []FileString
	for _,f := range dfs.Files {
		fmt.Println(f.Pinned)
		f2 := FileString{f.Filename, string(f.Data), true}
		r = append(r, f2)
	}
	return r
}
