package main

import (
	"fmt"
	"./kademlia"
	"github.com/ccding/go-stun/stun"
  "time"
	"os"
	"strconv"
	"math/rand"
)

func main() {

	args := os.Args[1:]
	port := os.Getenv("PORT")
	ip := "172.17.0.1"
	if port == "" {
		port = args[0]
		// Use STUN to get external IP address
		_, host, _ := stun.NewClient().Discover()
		ip = host.IP()
	}



	kID := kademlia.NewHashKademliaID(ip+":"+port)
	fmt.Printf("Starting node with id %v\n", kID.String())

	// Create Kademlia struct object
	me := kademlia.NewContact(kID, ip+":"+port)
  me.CalcDistance(kID)
	rt := kademlia.NewRoutingTable(me)

	dfs := kademlia.DFS{rt}
	network := kademlia.Network{ip, port}
	queue := kademlia.Queue{make(chan kademlia.Contact), rt}
	go queue.Run()

	k := kademlia.Kademlia{rt, &network, &queue, &dfs}
	// Starting listening
	go kademlia.HTTPListen(port, &k)
	go kademlia.Listen(&k)

  time.Sleep(time.Second * 1)

	if port != "8000" {

		contact_port_i,_ := strconv.Atoi(port)
		contact_port_i -= 8000
		contact_port_i = 8000+rand.Intn(contact_port_i/2)
		contact_port := strconv.Itoa(contact_port_i)

		c := kademlia.NewContact(kademlia.NewRandomKademliaID(), ip+":"+contact_port)

		fmt.Println("First contanct is: "+c.String())
		done := make(chan *kademlia.KademliaID)
		go k.Network.SendPingMessage(&me, &c, done)
		id := <- done
		c.ID = id
		queue.Enqueue(c)

		k.LookupContact(&me)
	}

	fmt.Println("Ready for use!")

	select {}
}
