package kademlia


import (
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
	Waiting 			chan File
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
		Waiting: make(chan File, 10000),
		Files: make(map[string]File),
	}
}

func (dfs *DFS) InitDFS(kademlia Kademlia) {
	dfs.Kademlia = kademlia

	go func ()  {
		for {
			select {
				case f, ok := <- dfs.Waiting:
			  	if ok {
						if _,ok := dfs.Files[f.Filename]; !ok {
							fmt.Printf("Writing %v to storage\n", f.Filename)
							dfs.Files[f.Filename] = f
							go dfs.PurgeFile(f.Filename)
						}
			  	}
			  default:
			}
			time.Sleep(time.Millisecond * 10)
		}
	}()
}

func (dfs *DFS) Cat(filename string, done chan []byte) {
	val, ok := dfs.Files[filename]
	if ok {
  	done<-val.Data
	} else {
		done<-[]byte("CatFileDoesntExists")
	}
}

func (dfs *DFS) Store(filename string, data []byte, done chan string) {
	// Lock so only one goroutine at a time can write.
	dfs.Waiting <- File{filename, data, make(chan bool)}
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

		val, ok := dfs.Files[filename]
		if ok {
			data := val.Data
			delete(dfs.Files, filename)

			target := NewContact(NewHashKademliaID(filename), "localhost:1111")
			contacts := dfs.Kademlia.FindClosestInCluster(&target)

			cut := 20

			if len(contacts) < cut {
				cut = len(contacts)
			}

			for _, c := range contacts[:cut] {
				doneStore := make(chan string)
				go dfs.Kademlia.Network.SendStoreMessage(&dfs.Kademlia.RoutingTable.me,
					&c, filename, data, doneStore)
				<-doneStore
			}
		}

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
