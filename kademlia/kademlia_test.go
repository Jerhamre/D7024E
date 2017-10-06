package kademlia

import (
  "fmt"
  "testing"
  "time"
)

func TestKademliaLookupContact (t *testing.T){
  fmt.Println("-----Test-LookupContact-----")
  ip := "localhost"

  // setup node 1
	port := "8100"
	kID := NewHashKademliaID(ip+":"+port)
	//fmt.Printf("Starting node with id %v\n", kID.String())
	// Create Kademlia struct object
	me := NewContact(kID, ip+":"+port)
  me.CalcDistance(kID)
	rt := NewRoutingTable(me)
	dfs := NewDFS(rt, 10)
	network := Network{ip, port}
	queue := Queue{make(chan Contact), rt, 10}
	go queue.Run()
	k := Kademlia{rt, &network, &queue, &dfs}
	dfs.InitDFS(k)


  // Setup node 2
  port2 := "8101"
  kID2 := NewHashKademliaID(ip+":"+port2)
  //fmt.Printf("Starting node with id %v\n", kID.String())
  // Create Kademlia struct object
	me2 := NewContact(kID2, ip+":"+port2)
  me2.CalcDistance(kID2)
	rt2 := NewRoutingTable(me2)
	dfs2 := NewDFS(rt2, 10)
	network2 := Network{ip, port2}
	queue2 := Queue{make(chan Contact), rt2, 10}
	go queue2.Run()
	k2 := Kademlia{rt2, &network2, &queue2, &dfs2}
	dfs2.InitDFS(k2)


  go Listen(&k)
  go Listen(&k2)

  time.Sleep(time.Second * 1)

  // add node 2 to node 1 RoutingTable
  queue.Enqueue(me2)

  // lookup should add node 1 to node 2s routingTable
  k.LookupContact(&me)
  fmt.Println("LOOKUPCONTACT OK!")
}

func TestKademliaLookupData (t *testing.T){
  fmt.Println("-----Test-LookupData-----")
  ip := "localhost"

  // setup node 1
  port := "8102"
  kID := NewHashKademliaID(ip+":"+port)
  //fmt.Printf("Starting node with id %v\n", kID.String())
  // Create Kademlia struct object
  me := NewContact(kID, ip+":"+port)
  me.CalcDistance(kID)
  rt := NewRoutingTable(me)
  dfs := NewDFS(rt, 10)
  network := Network{ip, port}
  queue := Queue{make(chan Contact), rt, 10}
  go queue.Run()
  k := Kademlia{rt, &network, &queue, &dfs}
  dfs.InitDFS(k)


  // Setup node 2
  port2 := "8103"
  kID2 := NewHashKademliaID(ip+":"+port2)
  fmt.Printf("Starting node with id %v\n", kID.String())
  // Create Kademlia struct object
  me2 := NewContact(kID2, ip+":"+port2)
  me2.CalcDistance(kID2)
  rt2 := NewRoutingTable(me2)
  dfs2 := NewDFS(rt2, 10)
  network2 := Network{ip, port2}
  queue2 := Queue{make(chan Contact), rt2, 10}
  go queue2.Run()
  k2 := Kademlia{rt2, &network2, &queue2, &dfs2}
  dfs2.InitDFS(k2)



  go Listen(&k)
  go Listen(&k2)

  time.Sleep(time.Second * 1)

  // add each node to eachothers RoutingTable
  queue.Enqueue(me2)
  queue2.Enqueue(me)

  done := make(chan string)
  testByte := []byte("file content")
  go k2.DFS.Store("memes", testByte, done)
  <-done

  doneLookup := make(chan []byte)
  go k.LookupData("memes", doneLookup)
  result := <-doneLookup
  //<-doneLookup
  //fmt.Println("result",result)
  if(string(result) != "file content"){
    t.Fatal("lookup unsuccessful",result)
  }
  fmt.Println("LOOKUPDATA OK!")
}

func TestKademliaFindClosestInCluster (t  *testing.T){
  // TODO
}

func TestKademliaStore (t *testing.T){
  fmt.Println("-----Test-Store-----")
  ip := "localhost"

  // setup node 1
  port := "8104"
  kID := NewHashKademliaID(ip+":"+port)
  //fmt.Printf("Starting node with id %v\n", kID.String())
  // Create Kademlia struct object
  me := NewContact(kID, ip+":"+port)
  me.CalcDistance(kID)
  rt := NewRoutingTable(me)
  dfs := NewDFS(rt, 10)
  network := Network{ip, port}
  queue := Queue{make(chan Contact), rt, 10}
  go queue.Run()
  k := Kademlia{rt, &network, &queue, &dfs}
  dfs.InitDFS(k)

  go Listen(&k)

  time.Sleep(time.Second * 1)

  storeDone := make(chan string)
  testData := []byte("file content")
  go k.Store("test", testData, storeDone)
  res := <-storeDone

  if res != "test" {
    t.Error("store failed")
  }
  fmt.Println("STORE OK!")
}

func TestKademliaPin (t *testing.T){
  fmt.Println("-----Test-Pin-----")
  ip := "localhost"

  // setup node 1
  port := "8107"
  kID := NewHashKademliaID(ip+":"+port)
  //fmt.Printf("Starting node with id %v\n", kID.String())
  // Create Kademlia struct object
  me := NewContact(kID, ip+":"+port)
  me.CalcDistance(kID)
  rt := NewRoutingTable(me)
  dfs := NewDFS(rt, 10)
  network := Network{ip, port}
  queue := Queue{make(chan Contact), rt, 10}
  go queue.Run()
  k := Kademlia{rt, &network, &queue, &dfs}
  dfs.InitDFS(k)

  go Listen(&k)

  time.Sleep(time.Second * 1)

  storeDone := make(chan string)
  testData := []byte("file content")
  go k.Store("pin", testData, storeDone)
  res := <-storeDone

  if res != "pin" {
    t.Fatal("store in TestKademliaPin failed")
  }
  pinDone := make(chan bool)

  go k.Pin("pin", pinDone)
  resPin := <-pinDone
  if !resPin {
    t.Fatal("pin failed")
  }
  fmt.Println("PIN OK!")

}

func TestKademliaUnpin (t *testing.T){
  fmt.Println("-----Test-Unpin-----")
  ip := "localhost"

  // setup node 1
  port := "8108"
  kID := NewHashKademliaID(ip+":"+port)
  //fmt.Printf("Starting node with id %v\n", kID.String())
  // Create Kademlia struct object
  me := NewContact(kID, ip+":"+port)
  me.CalcDistance(kID)
  rt := NewRoutingTable(me)
  dfs := NewDFS(rt, 10)
  network := Network{ip, port}
  queue := Queue{make(chan Contact), rt, 10}
  go queue.Run()
  k := Kademlia{rt, &network, &queue, &dfs}
  dfs.InitDFS(k)

  go Listen(&k)

  time.Sleep(time.Second * 1)

  storeDone := make(chan string)
  testData := []byte("file content")
  go k.Store("unpin", testData, storeDone)
  res := <-storeDone

  if res != "unpin" {
    t.Fatal("store in TestKademliaUnpin failed")
  }
  pinDone := make(chan bool)

  go k.Pin("unpin", pinDone)
  resPin := <-pinDone
  if !resPin {
    t.Fatal("pin in TestKademliaUnpin failed")
  }

  unpinDone := make(chan bool)
  go k.Unpin("unpin",unpinDone)
  resUnpin := <- unpinDone
  if !resUnpin {
    t.Fatal("unpin failed")
  }


  fmt.Println("UNPIN OK!")
}
