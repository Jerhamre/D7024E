package kademlia


import (
	"fmt"
	"sync"
  "io/ioutil"
	"os"
	"time"
)

type File struct {
	filename string
	content []byte
	pinned bool
}

type DFS struct {
	RoutingTable 	*RoutingTable
	PurgeTimer 		int
	Kademlia			Kademlia
	mux 					sync.Mutex // ← this mutex protects the cache below
	PurgeList			map[string]chan bool
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
		PurgeList: make(map[string]chan bool),
	}
}

func (dfs *DFS) InitDFS(kademlia Kademlia) {

	dfs.Kademlia = kademlia

	fmt.Println("Init DFS")

	os.Mkdir("files/"+dfs.RoutingTable.me.Address+"/", 0775)

	files, err := ioutil.ReadDir("files/"+dfs.RoutingTable.me.Address+"/")
  if err != nil {
    fmt.Println(err)
  }

  for _, f := range files {
		go dfs.PurgeFile(f.Name())
  }
}

func (dfs *DFS) Cat(hash string, done chan []byte) {

	if _, err := os.Stat("files/"+dfs.RoutingTable.me.Address+"/"+hash); os.IsNotExist(err) {
		done<-[]byte("CatFileDoesntExists")
		return
	}
	dat, err := ioutil.ReadFile("files/"+dfs.RoutingTable.me.Address+"/"+hash)
	if err != nil {
		 fmt.Printf("Cat %v\n", err)
		 done<-[]byte("CatError")
	}
  done<-dat
}

func (dfs *DFS) Store(filename string, data []byte, done chan string) {
	// Lock so only one goroutine at a time can write.
	dfs.mux.Lock()

	err := ioutil.WriteFile("files/"+dfs.RoutingTable.me.Address+"/"+filename, data, 0644)
	if err != nil {
		fmt.Printf("Store %v\n", err)
		done<-"StoreError"
	}

	go dfs.StopPurge(filename)
	go dfs.PurgeFile(filename)

	dfs.mux.Unlock()

	done<-filename
}

func (dfs *DFS) Pin(hash string, done chan bool) {
		go dfs.StopPurge(hash)
    done<-true
}

func (dfs *DFS) Unpin(hash string, done chan bool) {
	go dfs.PurgeFile(hash)
  done<-true
}

func (dfs *DFS) PurgeFile(filename string) {
	dfs.PurgeList[filename] = make(chan bool)
	select {
		case <-dfs.PurgeList[filename]:
			fmt.Println("Purge interrupted: "+filename)
		case <-time.After(time.Duration(dfs.PurgeTimer) * time.Second):


			// Get file content before Remove
			done := make(chan []byte)
			go dfs.Cat(filename, done)
			data := <-done

		  fmt.Println("File purged: "+filename)
			os.Remove("files/"+dfs.RoutingTable.me.Address+"/"+filename)

			fmt.Println("Republish file")
			target := NewContact(NewHashKademliaID(filename), "localhost:1111")
			fmt.Println(target.String())
			contacts := dfs.Kademlia.FindClosestInCluster(&target)

			for _, c := range contacts {
				doneStore := make(chan string)
				go dfs.Kademlia.Network.SendStoreMessage(&dfs.Kademlia.RoutingTable.me, &c, filename, data, doneStore)
				<-doneStore
			}
	}
}

func (dfs *DFS) StopPurge(filename string) {
	dfs.PurgeList[filename]<-true
}

func (dfs *DFS) GetFiles() {
	files, err := ioutil.ReadDir("./files/"+dfs.RoutingTable.me.Address+"/")
    if err != nil {
        fmt.Println(err)
    }

    for _, f := range files {
            fmt.Println(f.Name())
    }
}
