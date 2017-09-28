package main

import (
	"fmt"
	"./kademlia"
	//"github.com/ccding/go-stun/stun"
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

	dfs := kademlia.NewDFS(rt, 3)
	dfs.InitDFS()

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
		contact_port_i = 8000+rand.Intn(contact_port_i)
		contact_port := strconv.Itoa(contact_port_i)

		c := kademlia.NewContact(kademlia.NewRandomKademliaID(), ip+":"+contact_port)

		done := make(chan *kademlia.KademliaID)
		go k.Network.SendPingMessage(&me, &c, done)
		id := <- done
		c.ID = id
		fmt.Println("First contact is: "+c.String())
		queue.Enqueue(c)

		k.LookupContact(&me)


		time.Sleep(3*time.Second)

		fmt.Println("Store main.go")
		store := make(chan string)
		go dfs.Store("test", []byte("file content"), store)
		s := <-store
		fmt.Println(s)

		time.Sleep(2*time.Second)
		dfs.PurgeList["test"]<-true

		fmt.Println("Cat main.go")
		cat :=make(chan []byte)
		go dfs.Cat("test", cat)
		fmt.Println(string(<-cat))
	}

	fmt.Println("Ready for use!")

	select {}
}
