package kademlia


import (
	"fmt"
	"sync"
  "io/ioutil"
	"os"
	"time"
)

type DFS struct {
	RoutingTable 	*RoutingTable
	PurgeTimer 		int
	mux 					sync.Mutex // ← this mutex protects the cache below
	PurgeList			map[string]chan bool
}

func NewDFS(RoutingTable 	*RoutingTable, PurgeTimer int) DFS{
 return DFS{
	 RoutingTable: RoutingTable,
	 PurgeTimer: PurgeTimer,
	 PurgeList: make(map[string]chan bool),
 }
}
func check(e error) {
    if e != nil {
        panic(e)
    }
}

func (dfs *DFS) Cat(hash string, done chan []byte) {

	if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
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

	d1 := []byte("hello\ngo2\n")
	err := ioutil.WriteFile("files/"+dfs.RoutingTable.me.Address+"/"+filename, d1, 0644)
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

    // TODO

    done<-true
}

func (dfs *DFS) Unpin(hash string, done chan bool) {

  // TODO

  done<-true
}

func (dfs *DFS) InitDFS() {

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

func (dfs *DFS) PurgeFile(filename string) {
	dfs.PurgeList[filename] = make(chan bool)
	select {
		case m := <-dfs.PurgeList[filename]:
			m = m
			fmt.Println("purge of "+filename+" interrupted")
		case <-time.After(time.Duration(dfs.PurgeTimer) * time.Second):
		  fmt.Println(filename+" purged")
			os.Remove("files/"+dfs.RoutingTable.me.Address+"/"+filename)
	}
}

func (dfs *DFS) StopPurge(filename string) {
	dfs.PurgeList[filename]<-true
}
