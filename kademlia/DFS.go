package kademlia


import (
	"fmt"
	"sync"
  "io/ioutil"
	"os"
	"time"
)

type File struct {
	Filename string
	Content string
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

	os.Mkdir("files/", 0775)
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
			break
		case <-time.After(time.Duration(dfs.PurgeTimer) * time.Second):


			// Get file content before Remove
			done := make(chan []byte)
			go dfs.Cat(filename, done)
			data := <-done

		  os.Remove("files/"+dfs.RoutingTable.me.Address+"/"+filename)

			target := NewContact(NewHashKademliaID(filename), "localhost:1111")
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

func (dfs *DFS) GetFiles() []File {
	var retFiles []File

	files, err := ioutil.ReadDir("./files/"+dfs.RoutingTable.me.Address+"/")
    if err != nil {
        fmt.Println(err)
    }

    for _, f := range files {
			filename := f.Name()
			done := make(chan []byte)
			go dfs.Cat(filename, done)
			content := <- done
			retFiles = append(retFiles, File{filename, string(content)})
    }

		return retFiles
}
