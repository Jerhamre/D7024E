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

	//var contacts []Contact
	var seenIDs []*KademliaID

	done := make(chan []Contact)

  //pick out alpha (3) nodes from its closest non-empty k-bucket
	returned_contacts := kademlia.RoutingTable.FindClosestContacts(target.ID, alpha)


  // send parallel, async FIND_NODE to the alpha (3) nodes chosen
	findNode := 0	// amount of current findNode calls
	findNodeResponses := 0 // recieved responses

  for _, contact := range returned_contacts {
		findNode++
		seenIDs = append(seenIDs, contact.ID)
		fmt.Println("ID", contact.Address)
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
		findNodeResponses++
		findNode--

		for _, node := range nodes {
			// checks if node has already been seen, if so, ignore
			if !contains(seenIDs, node.ID) {
				kademlia.RoutingTable.AddContact(node)
				returned_contacts = append(returned_contacts, node)
				//seenIDs = append(seenIDs, node.ID)
			}
		}

		// resends FIND_NODE to nodes that it has learned from previous RPC
		for findNode < alpha && findNodeResponses < k{

			// sends new RPC to the closest contact in the priority queue that has not been seen
			contacts := kademlia.RoutingTable.FindClosestContacts(target.ID, k)
			// TODO: find closest contact from returned_contacts NOT the routingTable
			// CalcDistance to target, Sort
			//contacts := returned_contacts.Sort()

			// creates async FIND_NODE to the first unseen contact in the routingTable
			for _, contact := range contacts{
				if !contains(seenIDs, contact.ID){
					findNode++
					go kademlia.Network.SendFindContactMessage(&kademlia.RoutingTable.me, &contact, target, done)
					break
				}
			}

		}

	}

}

// helper function for finding out if an element exists in a slice
// from stackoverflow https://stackoverflow.com/questions/10485743/contains-method-for-a-slice
// modified to take *KademliaID instead of string
func contains(slice []*KademliaID, item *KademliaID) bool {
    set := make(map[*KademliaID]struct{}, len(slice))
    for _, s := range slice {
        set[s] = struct{}{}
    }

    _, ok := set[item]
    return ok
}

func (kademlia *Kademlia) LookupData(hash string, done chan []byte) {
	kademliaID := NewHashKademliaID(hash)
	kademliaID = kademliaID

	// TODO

	done<-[]byte("file content")
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
