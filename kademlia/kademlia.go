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
	returned_contacts := kademlia.FindClosestInCluster(target)

	for _, contact := range returned_contacts {
		kademlia.Queue.Enqueue(contact)
	}

}

func (kademlia *Kademlia) LookupData(hash string, done chan []byte) {
	//fmt.Println("hash",hash)
	kademliaID := NewHashKademliaID(hash)
	kademliaID = kademliaID

	fmt.Println("LookupData", hash)

	const alpha = 3
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
	go kademlia.DFS.Cat(hash, dfsCat)
	stored := <-dfsCat
	fmt.Println("DFS CAT RETURNED:",string(stored))
	if string(stored) != "CatError" && string(stored) != "CatFileDoesntExists" && string(stored) != ""{
		// file exists on node, done
		fmt.Println("file found on node, terminating lookupData: ", stored)
		done<-stored
		return
	}
	fmt.Println("file not found on local, contacting cluster")
	returned_contacts := ContactCandidates{kademlia.RoutingTable.FindClosestContacts(kademliaID, k)}


	//fmt.Println("Searching in cluster for: ",kademliaID)
	noMoreQueable := false
	for {
		for active == alpha || i == k || noMoreQueable{
			temp := <-ch
			noMoreQueable = false
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
				//done<-"data not found"
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
				} else if active < alpha{
					noMoreQueable = true
				}
			}
			//fmt.Println("active",active)
			if active == 0{
				fmt.Println("lookupData terminated, Data not found")
				//done<-"data not found"
				return
			}
		}
	}

	//done<-lookupRes
}

func (kademlia *Kademlia) FindClosestInCluster(target *Contact) []Contact{
	//fmt.Println("FindClosestInCluster running")
	const alpha = 3
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
	for _, contact := range returned_contacts.contacts {
		seen[contact.ID.String()] = struct{}{}
	}
	noMoreQueable := false
	for {
		for active == alpha || i == k || noMoreQueable{
			temp := <-ch
			noMoreQueable = false
			//fmt.Println("SFCM returned", temp)

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
		returned_contacts.Sort()

		if i < k {

			for _, contact := range returned_contacts.contacts {
				if _, ok := sent[contact.ID.String()]; !ok{
					fmt.Println("sending SFCM to", contact.Address)
					go kademlia.Network.SendFindContactMessage(&kademlia.RoutingTable.me, &contact, target, ch)
					sent[contact.ID.String()] = struct{}{}
					active++
					i++
					break
				} else if active < alpha{
					noMoreQueable = true
				}
			}
			if active == 0{
				//fmt.Println("FindClosestInCluster terminated")
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
	fmt.Println("Storing file ",filename)

	go kademlia.DFS.Store(filename, data, localStoreDone)
	store := <-localStoreDone

	// Find the closest nodes in the cluster
	dataContact := NewContact(kademliaID, "data")
	closestNodes := kademlia.FindClosestInCluster(&dataContact)
	SSMSent := 0
	//fmt.Println("closestNodes",closestNodes)
	for _, c := range closestNodes{
		fmt.Println("SSM sent to", c.Address)
		go kademlia.Network.SendStoreMessage(&kademlia.RoutingTable.me, &c, filename, data, SSM)
		SSMSent++
	}
	sentSSM := 0

	for sentSSM < SSMSent {
		response := <-SSM
		sentSSM++
		fmt.Println(response)
		if response != "fail"{
			// atleast one node has the file
			fileStored = true
			store = response
			fmt.Println("Store successful on one node in cluster")
			break
		}

	}
	if !fileStored {
		fmt.Println("only stored locally, all SSM failed")
	}
	fmt.Println("store terminated")
	done<-store
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
