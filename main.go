package main

import (
	"fmt"
	"./kademlia"
	"github.com/ccding/go-stun/stun"
  "time"
	"os"
)

func main() {

	ports := os.Args[1:]
	portKademlia := ports[0]
	portHTTP := ports[1]

	kID := kademlia.NewRandomKademliaID()
	fmt.Printf("Starting node with id %v\n", kID.String())

	// Use STUN to get external IP address
	_, host, _ := stun.NewClient().Discover()
	ip := host.IP()

	// Create Kademlia struct object
	me := kademlia.NewContact(kID, ip+":"+portKademlia)
  me.CalcDistance(kID)
	rt := kademlia.NewRoutingTable(me)

  /*
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
  testContact := kademlia.NewContact(kademlia.NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8003")
	rt.AddContact(testContact)
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8004"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8005"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8006"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8007"))
	*/

	dfs := kademlia.DFS{rt}
	network := kademlia.Network{ip, portKademlia}
	queue := kademlia.Queue{make(chan kademlia.Contact), rt}
	go queue.Run()

	k := kademlia.Kademlia{rt, &network, &queue, &dfs}
	// Starting listening
	go kademlia.HTTPListen(portHTTP, &k)
	go kademlia.Listen(&k)

  time.Sleep(time.Second * 1)

	if portKademlia != "8000" {
		c := kademlia.NewContact(kademlia.NewRandomKademliaID(), ip+":8000")
		done := make(chan *kademlia.KademliaID)
	  go k.Network.SendPingMessage(&me, &c, done)
		for {
			id := <- done
			c.ID = id
			queue.Enqueue(c)

			fmt.Println("SendFindContactMessage")
			done := make(chan []kademlia.Contact)
			go network.SendFindContactMessage(&me, &c, &me, done)
			fmt.Println(<-done)

			fmt.Println("SendStoreMessage")
			done2 := make(chan string)
			go network.SendStoreMessage(&me, &c, "qwe", []byte("file content"), done2)
			fmt.Println(<-done2)

			fmt.Println("SendFindDataMessage")
			done3 := make(chan []byte)
			go network.SendFindDataMessage(&me, &c, "qwe", done3)
			fmt.Println(string(<-done3))

			//k.LookupContact(&me)
		}
	}

	fmt.Println("Ready for use!")

	select {}
}
