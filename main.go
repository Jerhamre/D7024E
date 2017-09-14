package main

import (
	"fmt"
	"./kademlia"
	"github.com/ccding/go-stun/stun"
  "time"
)

func main() {
	kID := kademlia.NewRandomKademliaID()
	fmt.Printf("Started node with id %v\n", kID.String())

	// Use STUN to get external IP address
	_, host, _ := stun.NewClient().Discover()
	ip := host.IP()
	port := "8001"

	// Create Kademlia struct object
	contact := kademlia.NewContact(kID, ip+":"+port)
  contact.CalcDistance(kID)
	rt := kademlia.NewRoutingTable(contact)

  /*rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
  testContact := kademlia.NewContact(kademlia.NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8003")
	rt.AddContact(testContact)
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8004"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8005"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8006"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8007"))
*/
	network := kademlia.Network{ip, port}
	k := kademlia.Kademlia{rt, &network}

	// Starting listening
	go kademlia.HTTPListen("8081", &k)
	go kademlia.Listen(&k)

  time.Sleep(time.Second * 1)


	c := kademlia.NewContact(kademlia.NewRandomKademliaID(), ip+":8000")
	rt.AddContact(c)
	k.LookupContact(&contact)

	/*
  //k.Network.SendPingMessage(&c)
  done := make(chan []kademlia.Contact)
  go k.Network.SendFindContactMessage(&c, &testContact, done)

  for {
    fmt.Println("done")
    conts := <- done
    fmt.Println("asdf"+conts[0].String())
  }*/
	for {}
}
