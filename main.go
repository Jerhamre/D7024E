package main

import (
	"fmt"
	"./kademlia"
	//"github.com/ccding/go-stun/stun"
  "time"
	"os"
)

func main() {

	args := os.Args[1:]
	port := os.Getenv("PORT")
	ip := "172.17.0.1"
	if port == "" {
		port = args[0]
		// Use STUN to get external IP address
		//_, host, _ := stun.NewClient().Discover()
		//ip = host.IP()
		ip = "localhost"
	}

	kID := kademlia.NewHashKademliaID(ip+":"+port)
	fmt.Printf("Starting node with id %v\n", kID.String())

	// Create Kademlia struct object
	me := kademlia.NewContact(kID, ip+":"+port)
  me.CalcDistance(kID)
	rt := kademlia.NewRoutingTable(me)

	dfs := kademlia.NewDFS(rt, 10)

	network := kademlia.Network{ip, port}
	queue := kademlia.Queue{make(chan kademlia.Contact), rt}
	go queue.Run()

	k := kademlia.Kademlia{rt, &network, &queue, &dfs}

	dfs.InitDFS(k)

	// Starting listening
	go kademlia.HTTPListen(port, &k)
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

		k.LookupContact(&me)


		time.Sleep(1*time.Second)
		if port == "8001" {
			fmt.Println("Store main.go")
			store := make(chan string)
			go dfs.Store("test", []byte("file content"), store)
			s := <-store
			fmt.Println(s)
		}

		/*

		time.Sleep(2*time.Second)

		fmt.Println("Pin main.go")
		pin := make(chan bool)
		go dfs.Pin("test", pin)
		fmt.Println(<-pin)


		time.Sleep(4*time.Second)

		fmt.Println("Unpin main.go")
		unpin := make(chan bool)
		go dfs.Unpin("test", unpin)
		fmt.Println(<-unpin)

		*/
	}

	fmt.Println("Ready for use!")

	select {}
}
