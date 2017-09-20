package kademlia


import (
)

type DFS struct {
	RoutingTable 	*RoutingTable
}

func (dfs *DFS) Cat(hash string, done chan []byte) {

  // TODO

  done<-[]byte("file content")
}

func (dfs *DFS) Store(filename string, data []byte, done chan string) {

  // TODO

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
