package main

import (
	"fmt"
	"d7024e/kademlia"
  "time"
	"os"
)

func main() {

	args := os.Args[1:]
	port := os.Getenv("PORT")
	ip := "172.17.0.1"
	if port == "" {
		port = args[0]
		ip = "localhost"
	}

	kID := kademlia.NewHashKademliaID(ip+":"+port)
	fmt.Printf("Starting node with id %v\n", kID.String())

	// Create Kademlia struct object
	me := kademlia.NewContact(kID, ip+":"+port)
  me.CalcDistance(kID)
	rt := kademlia.NewRoutingTable(me)

	dfs := kademlia.NewDFS(rt, 10000)

	network := kademlia.Network{ip, port}
	queue := kademlia.Queue{make(chan kademlia.Contact, 10000), rt, 10}
	go queue.Run()

	k := kademlia.Kademlia{rt, &network, &queue, &dfs}

	dfs.InitDFS(k)

	// Starting listening
	go kademlia.HTTPListen(port, &k, "kademlia/templates")
	go kademlia.Listen(&k)

  time.Sleep(time.Second * 1)

	if port != "8000" {
		c := kademlia.NewContact(kademlia.NewRandomKademliaID(), ip+":8000")
		done := make(chan *kademlia.KademliaID)
		go k.Network.SendPingMessage(&me, &c, done)
		id := <- done
		c.ID = id
		fmt.Println("First contact is: "+c.String())
		queue.Enqueue(c)

	  time.Sleep(time.Second * 1)

		k.LookupContact(&me)

		if port == "8001" {
			done2 := make(chan string)
			dfs.Store("qweqwe", []byte("content DLC"), done2)
			fmt.Println(<-done2)
		}
	}

	fmt.Println("Ready for use!")

	select {}
}
