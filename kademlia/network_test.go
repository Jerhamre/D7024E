package kademlia

import (
	"testing"
	"time"
)

func TestSendPingMessage(t *testing.T) {
	port1 := "8000"
	port2 := "8001"
	k := HelperNetworkListen(port1, port2)
	c := NewContact(NewRandomKademliaID(), "localhost:"+port2)
	done := make(chan *KademliaID)
	go k.Network.SendPingMessage(&k.RoutingTable.me, &c, done)
	id := <- done
	c.ID = id

	expectedID := "6c6f63616c686f73743a38303031202020202020"

	if(c.ID.String() != expectedID) {
    t.Fatalf("Expected %v but got %v", expectedID, c.ID.String())
  }
}

func TestSendFindContactMessage(t *testing.T) {
	port1 := "8002"
	port2 := "8003"
	k := HelperNetworkListen(port1, port2)
	c := NewContact(NewRandomKademliaID(), "localhost:"+port2)
	done := make(chan []Contact)
	go k.Network.SendFindContactMessage(&k.RoutingTable.me, &c, &k.RoutingTable.me, done)
	contacts := <- done

	var expected []Contact
	expected = append(expected, k.RoutingTable.me)

	if(contacts[0].String() != k.RoutingTable.me.String()) {
    t.Fatalf("Expected %v but got %v", k.RoutingTable.me.String(), contacts[0].String())
  }
}

func TestSendStoreMessage(t *testing.T) {
	port1 := "8004"
	port2 := "8005"
	k := HelperNetworkListen(port1, port2)
	c := NewContact(NewRandomKademliaID(), "localhost:"+port2)

	filename := "test"
	data := []byte("file content555")

	done := make(chan string)
	go k.Network.SendStoreMessage(&k.RoutingTable.me, &c, filename, data, done)
	ret := <- done

	if(ret != filename) {
    t.Fatalf("Expected %v but got %v", filename, ret)
  }
}

func TestSendFindDataMessage(t *testing.T) {
	port1 := "8006"
	port2 := "8007"
	k := HelperNetworkListen(port1, port2)
	c := NewContact(NewRandomKademliaID(), "localhost:"+port2)

	filename := "test"
	data := []byte("file content666")

	done := make(chan string)
	go k.Network.SendStoreMessage(&k.RoutingTable.me, &c, filename, data, done)
	ret := <- done

	time.Sleep(time.Millisecond * 30)

	if(ret != filename) {
    t.Fatalf("Expected Store %v but got %v", filename, ret)
  }

	done2 := make(chan []byte)
	go k.Network.SendFindDataMessage(&k.RoutingTable.me, &c, filename, done2)
	ret2 := <-done2

	if(!testEq(ret2, data)) {
    t.Fatalf("Expected Find %v but got %v", string(data), string(ret2))
  }
}

func TestSendPinMessage(t *testing.T) {
	port1 := "8008"
	port2 := "8009"
	k := HelperNetworkListen(port1, port2)
	c := NewContact(NewRandomKademliaID(), "localhost:"+port2)

	filename := "test"
	data := []byte("file content")

	done := make(chan string)
	go k.Network.SendStoreMessage(&k.RoutingTable.me, &c, filename, data, done)
	ret := <- done

	if(ret != filename) {
    t.Fatalf("Expected Store %v but got %v", filename, ret)
  }

	done2 := make(chan bool)
	go k.Network.SendPinMessage(&k.RoutingTable.me, &c, filename, done2)
	ret2 := <- done2

	if(ret2 != true) {
    t.Fatalf("Expected Pin %v but got %v", true, ret2)
  }
}

func TestSendUnpinMessage(t *testing.T) {
	port1 := "8010"
	port2 := "8011"
	k := HelperNetworkListen(port1, port2)
	c := NewContact(NewRandomKademliaID(), "localhost:"+port2)

	filename := "test"
	data := []byte("file content")

	done := make(chan string)
	go k.Network.SendStoreMessage(&k.RoutingTable.me, &c, filename, data, done)
	ret := <- done

	if(ret != filename) {
    t.Fatalf("Expected Store %v but got %v", filename, ret)
  }

	done2 := make(chan bool)
	go k.Network.SendPinMessage(&k.RoutingTable.me, &c, filename, done2)
	ret2 := <- done2

	if(ret2 != true) {
    t.Fatalf("Expected Pin %v but got %v", true, ret2)
  }

	done3 := make(chan bool)
	go k.Network.SendUnpinMessage(&k.RoutingTable.me, &c, filename, done3)
	ret3 := <- done3

	if(ret3 != true) {
    t.Fatalf("Expected Unpin %v but got %v", true, ret3)
  }
}

func TestHandleRequest(t *testing.T) {

}

func HelperNetworkListen(port1 string, port2 string) *Kademlia{
	kID1 := NewHashKademliaID("localhost:"+port1)
  me1 := NewContact(kID1, "localhost:"+port1)
  me1.CalcDistance(kID1)
  rt1 := NewRoutingTable(me1)
  dfs1 := NewDFS(rt1, 10000)
  network1 := Network{"localhost", port1}
  queue1 := Queue{make(chan Contact), rt1, 10}
  k1 := Kademlia{rt1, &network1, &queue1, &dfs1}
  go queue1.Run()
	dfs1.InitDFS(k1)
	go Listen(&k1)

	kID2 := NewHashKademliaID("localhost:"+port2)
  me2 := NewContact(kID2, "localhost:"+port2)
  me2.CalcDistance(kID2)
  rt2 := NewRoutingTable(me2)
  dfs2 := NewDFS(rt2, 10000)
  network2 := Network{"localhost", port2}
  queue2 := Queue{make(chan Contact), rt2, 10}
  k2 := Kademlia{rt2, &network2, &queue2, &dfs2}
  go queue2.Run()
	dfs2.InitDFS(k2)
	go Listen(&k2)


	time.Sleep(time.Millisecond * 25)

	return &k1
}
