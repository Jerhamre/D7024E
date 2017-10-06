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
	val, ok := dfs.Files[filename]
	if ok {
  	done<-val.Data
	} else {
		done<-nil
	}
}

func (dfs *DFS) Store(filename string, data []byte, done chan string) {
	// Lock so only one goroutine at a time can write.
	dfs.mux.Lock()

	if _, ok := dfs.Files[filename]; !ok {
		dfs.Files[filename] = File{filename, data, make(chan bool)}

		//go dfs.StopPurge(filename)
		go dfs.PurgeFile(filename)
	}

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

		dfs.mux.Lock()

		val, ok := dfs.Files[filename]
		if ok {
			data := val.Data
			fmt.Println(dfs.Files[filename])
			delete(dfs.Files, filename)

			target := NewContact(NewHashKademliaID(filename), "localhost:1111")
			contacts := dfs.Kademlia.FindClosestInCluster(&target)

			fmt.Printf("Send file to %v contacts\n", len(contacts))
			fmt.Println(contacts)

			for _, c := range contacts[:20] {
				doneStore := make(chan string)
				go dfs.Kademlia.Network.SendStoreMessage(&dfs.Kademlia.RoutingTable.me,
					&c, filename, data, doneStore)
				<-doneStore
			}
		}

		dfs.mux.Unlock()

		dfs.Files[filename].Pinned<-false
	}
}

func (dfs *DFS) StopPurge(filename string) {
	dfs.Files[filename].Pinned <- true
}

func (dfs *DFS) GetFiles() []FileString {
	var r []FileString
	for _,f := range dfs.Files {
		f2 := FileString{f.Filename, string(f.Data), true}
		r = append(r, f2)
	}
	return r
}
