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
	const alpha = 3
	const k = 20

	seen := make(map[string]struct{})
	done := make(chan []Contact)

  //pick out alpha (3) nodes from its closest non-empty k-bucket
	returned_contacts := kademlia.RoutingTable.FindClosestContacts(target.ID, alpha)
	//fmt.Println("returned_contacts from RoutingTable", returned_contacts)

  // send parallel, async FIND_NODE to the alpha (3) nodes chosen
	findNode := 0	// amount of current findNode calls
	findNodeResponses := 0 // recieved responses

	seen[kademlia.RoutingTable.me.ID.String()] = struct{}{}
	//fmt.Println("me.ID",kademlia.RoutingTable.me.ID)

  for _, contact := range returned_contacts {
		findNode++
		seen[contact.ID.String()] = struct{}{}
		//fmt.Println("ID", contact.ID)
		//fmt.Println("Address", contact.Address)
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
			//fmt.Println("nodes return", nodes)
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
				//fmt.Println("enqueue on", node)
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
				//fmt.Println("sending new sendFindContactMessage to", contact)
				seen[contact.ID.String()] = struct{}{}
				go kademlia.Network.SendFindContactMessage(&kademlia.RoutingTable.me, &contact, target, done)
				break
			} else {
				//fmt.Println("seen (SFCM already sent to)", contact) // debugging contains
			}
		}

	}
	fmt.Println("LookupContact derminated")

}

func (kademlia *Kademlia) LookupData(hash string, done chan []byte) {
	//fmt.Println("hash",hash)
	kademliaID := NewHashKademliaID(hash)
	kademliaID = kademliaID

	fmt.Println("LookupData", hash)

	const alpha = 1
	const k = 20
	var active = 0
	ch := make(chan []Contact)
	dfsCat := make(chan []byte)
	SFDMresponse := make(chan []byte)
	seen := make(map[string]struct{})
	sent := make(map[string]struct{})

	i := 0
  goDone := 0
	seen[kademlia.RoutingTable.me.ID.String()] = struct{}{}
	dataContact := NewContact(kademliaID, "data")
	fmt.Println("cat on node 1")
	go kademlia.DFS.Cat(hash, dfsCat)
	fmt.Println("cat on node 2")
	stored := <-dfsCat
	if string(stored) != "CatError" && string(stored) != "CatFileDoesntExists"{
		// file exists on node, done
		fmt.Println("file found on node, terminating lookupData: ", stored)
		done<-stored
		return
	}

	returned_contacts := ContactCandidates{kademlia.RoutingTable.FindClosestContacts(kademliaID, k)}

	fmt.Println("lookupData")
	//fmt.Println("Searching in cluster for: ",kademliaID)

	for {
		for active == alpha || i == k {
			temp := <-ch
			//fmt.Println("ch return")

			// only add a contact if it has not been added before
			for _, node := range temp {
				if _, ok := seen[node.ID.String()]; !ok{
						returned_contacts.contacts = append(returned_contacts.contacts, node)
						seen[node.ID.String()] = struct{}{}
					}
				}

			active--
			goDone++
			if goDone == k{
				fmt.Println("lookupData terminated, data not found")
				return
			}
		}


		if i < k {

			for _, contact := range returned_contacts.contacts {
				if _, ok := sent[contact.ID.String()]; !ok{

					// sends a lookup data RPC to contact
					//fmt.Println("SFDM on", hash)
					go kademlia.Network.SendFindDataMessage(&kademlia.RoutingTable.me, &contact, hash, SFDMresponse)
					lookupRes := <-SFDMresponse

					// if file is found, return file and terminate
					if string(lookupRes) != "CatError" && string(lookupRes) != "CatFileDoesntExists"{
						fmt.Println("lookup Data successful, Terminating now: ", string(lookupRes))
						done<-lookupRes
						return
					}

					//fmt.Println("sending SFCM to", contact)
					go kademlia.Network.SendFindContactMessage(&kademlia.RoutingTable.me, &contact, &dataContact, ch)
					sent[contact.ID.String()] = struct{}{}
					active++
					i++
					break
				}
			}
			//fmt.Println("active",active)
			if active == 0{
				fmt.Println("lookupData terminated, Data not found")
				return
			}
		}
	}

	//done<-lookupRes
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

	for {
		for active == alpha || i == k {
			temp := <-ch
			//fmt.Println("ch return")

			// only add a contact if it has not been added before
			for _, node := range temp {
				if _, ok := seen[node.ID.String()]; !ok{
						returned_contacts.contacts = append(returned_contacts.contacts, node)
						seen[node.ID.String()] = struct{}{}
					}
				}

			active--
			goDone++
			if goDone == k{
				fmt.Println("FindClosestInCluster terminated")
				//fmt.Println("returned_contacts",returned_contacts.contacts)
				return returned_contacts.contacts
			}
		}


		if i < k {

			for _, contact := range returned_contacts.contacts {
				if _, ok := sent[contact.ID.String()]; !ok{
					//fmt.Println("sending SFCM to", contact)
					go kademlia.Network.SendFindContactMessage(&kademlia.RoutingTable.me, &contact, target, ch)
					sent[contact.ID.String()] = struct{}{}
					active++
					i++
					break
				}
			}
			//fmt.Println("active",active)
			if active == 0{
				//fmt.Println("FindClosestInCluster terminated")
				//fmt.Println("returned_contacts",returned_contacts.contacts)
				return returned_contacts.contacts
			}
		}
	}
}

func (kademlia *Kademlia) Store(filename string, data []byte, done chan string) {
	kademliaID := NewHashKademliaID(filename)
	var fileStored = false
	SSM := make(chan string, 20)
	localStoreDone := make(chan string)
	fmt.Println("Store",filename)

	go kademlia.DFS.Store(filename, data, localStoreDone)
	<-localStoreDone

	var returnedValue string
	// Find the closest nodes in the cluster
	dataContact := NewContact(kademliaID, "data")
	closestNodes := kademlia.FindClosestInCluster(&dataContact)
	SSMSent := 0

	for _, c := range closestNodes{
		go kademlia.Network.SendStoreMessage(&kademlia.RoutingTable.me, &c, filename, data, SSM)
		SSMSent++
	}
	sentSSM := 0

	for sentSSM < SSMSent {
		response := <-SSM
		sentSSM++
		if response != "fail"{
			// atleast one node has the file
			fileStored = true
			returnedValue = response
			fmt.Println("Store successful")
			break
		}

	}
	//fmt.Println("file stored:",fileStored)
	fmt.Println("Stored as: ",kademliaID)
	if !fileStored {
		// file has not been stored on any node
	}

	done<-returnedValue
}

func (kademlia *Kademlia) Pin(filename string, done chan bool) {
	kademliaID := NewHashKademliaID(filename)
	kademliaID = kademliaID
	cpin := make(chan bool)

	go kademlia.DFS.Pin(filename, cpin)
	response := <-cpin

	done<-response
}

func (kademlia *Kademlia) Unpin(filename string, done chan bool) {
	kademliaID := NewHashKademliaID(filename)
	kademliaID = kademliaID
	cunpin := make(chan bool)

	go kademlia.DFS.Unpin(filename, cunpin)
	response := <-cunpin

	done<-response
}
