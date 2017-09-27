package kademlia

import (
	"fmt"
)

type Kademlia struct {
	RoutingTable 	*RoutingTable
	Network 			*Network
	Queue 				*Queue
	DFS						*DFS
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	//returned_contacts = kademlia.FindClosestInCluster(target)


	const alpha = 3
	const k = 20

	seen := make(map[string]struct{})
	done := make(chan []Contact)

  //pick out alpha (3) nodes from its closest non-empty k-bucket
	returned_contacts := kademlia.RoutingTable.FindClosestContacts(target.ID, alpha)


  // send parallel, async FIND_NODE to the alpha (3) nodes chosen
	findNode := 0	// amount of current findNode calls
	findNodeResponses := 0 // recieved responses

	seen[kademlia.RoutingTable.me.ID.String()] = struct{}{}
	fmt.Println("me.ID",kademlia.RoutingTable.me.ID)

  for _, contact := range returned_contacts {
		findNode++
		seen[contact.ID.String()] = struct{}{}
		//fmt.Println("ID", contact.ID)
		fmt.Println("Address", contact.Address)
    go kademlia.Network.SendFindContactMessage(&kademlia.RoutingTable.me, &contact, target, done)
  }

	// if a round of FIND_NODES fails to return a node any closer than the closest
	// already seen, the initiator resends the FIND_NODE to all k-closest nodes
	// it has not already queried

	// The lookup terminates when the initiator has queried and gotten responses
	// from the k closest nodes it has seen.

	for findNode > 0 && findNodeResponses < k{
		nodes := <-done
		if len(nodes)!=0{
			fmt.Println("nodes return", nodes)
			findNodeResponses++
			findNode--
		} else {
			fmt.Println("SFCM failed")
			findNode--
		}


		for _, node := range nodes {
			// checks if node has already been seen, if so, ignore
			if _, ok := seen[node.ID.String()]; !ok{
				kademlia.Queue.Enqueue(node)
				fmt.Println("enqueue on", node)
				returned_contacts = append(returned_contacts, node)
			}
		}

		// resends FIND_NODE to nodes that it has learned from previous RPC
		// sends new RPC to the closest contact in the priority queue that has not been seen
		contacts := kademlia.RoutingTable.FindClosestContacts(target.ID, k)
		//sort.Sort(returned_contacts)
		// TODO: find closest contact from returned_contacts NOT the routingTable
		// CalcDistance to target, Sort
		//contacts := returned_contacts.Sort()

		// creates async FIND_NODE to the first unseen contact in the routingTable
		for _, contact := range contacts{
			if _, ok := seen[contact.ID.String()]; !ok{
				findNode++
				fmt.Println("sending new sendFindContactMessage to", contact)
				seen[contact.ID.String()] = struct{}{}
				//done := make(chan []Contact)
				go kademlia.Network.SendFindContactMessage(&kademlia.RoutingTable.me, &contact, target, done)
				break
			} else {
				fmt.Println("seen (SFCM already sent to)", contact) // debugging contains
			}
		}

	}
	fmt.Println("LookupContact derminated")

}

func (kademlia *Kademlia) LookupData(hash string, done chan []byte) {

	kademliaID := NewHashKademliaID(hash)
	kademliaID = kademliaID

	// Find the closest nodes in the cluster
	// closestNodes := FindClosestInCluster

	//for _, contact := range closestNodes{
		//go kademlia.Network.SendFindDataMessage(&kademlia.RoutingTable.me, &contact, filename string, done chan []byte)
	//}
	// TODO

	done<-[]byte("file content")
}

func (kademlia *Kademlia) FindClosestInCluster(target *Contact) []Contact{

	const alpha = 1
	const k = 20
	//var returned_SCFM = 0
	var active = 0
	//ch := make(chan []Contact, 3)
	ch := make(chan []Contact)
	seen := make(map[string]struct{})
	sent := make(map[string]struct{})

	i := 0
  goDone := 0
	seen[kademlia.RoutingTable.me.ID.String()] = struct{}{}

	returned_contacts := ContactCandidates{kademlia.RoutingTable.FindClosestContacts(target.ID, k)}
	//returned_contacts.contacts = append(returned_contacts.contacts, nodes ...)

	for {
		for active == alpha || i == k {
			temp := <-ch
			fmt.Println("ch return")

			// only add a contact if it has not been added before
			for _, node := range temp {
				if _, ok := seen[node.ID.String()]; !ok{
						returned_contacts.contacts = append(returned_contacts.contacts, node)
						seen[node.ID.String()] = struct{}{}
					}
				}
			//fmt.Println("contacts", returned_contacts.contacts)
			//returned_contacts.Sort()
			//fmt.Printf("ch: %v\n", temp)
			active--
			goDone++
			if goDone == k{
				fmt.Println("FindClosestInCluster terminated")
				fmt.Println("returned_contacts",returned_contacts.contacts)
				return returned_contacts.contacts
			}
		}


		if i < k {

			for _, contact := range returned_contacts.contacts {
				if _, ok := sent[contact.ID.String()]; !ok{
					fmt.Println("sending SFCM to", contact)
					go kademlia.Network.SendFindContactMessage(&kademlia.RoutingTable.me, &contact, target, ch)
					sent[contact.ID.String()] = struct{}{}
					active++
					i++
					break
				}
			}
			fmt.Println("active",active)
			if active == 0{
				fmt.Println("FindClosestInCluster terminated")
				fmt.Println("returned_contacts",returned_contacts.contacts)
				return returned_contacts.contacts
			}
		}
	}



	/*
	for returned_SCFM < 20 {
		for active == 3 {
			contacts<-ch
			returned_SCFM++
			active--
			returned_contacts = append(returned_contacts, contacts ...)
		}


		for _, contact = range returned_contacts && active != 3{
			active++
			go kademlia.Network.SendFindContactMessage(&kademlia.RoutingTable.me, &contact, target, ch)
		}

	}
	*/
}

func (kademlia *Kademlia) Store(filename string, data []byte, done chan string) {
	kademliaID := NewHashKademliaID(filename)

	// TODO

	done<-kademliaID.String()
}

func (kademlia *Kademlia) Pin(filename string, done chan bool) {
	kademliaID := NewHashKademliaID(filename)
	kademliaID = kademliaID

	// TODO

	done<-true
}

func (kademlia *Kademlia) Unpin(filename string, done chan bool) {
	kademliaID := NewHashKademliaID(filename)
	kademliaID = kademliaID

	// TODO

	done<-false
}
