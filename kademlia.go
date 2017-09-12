package d7024e

import (
	"container/heap"
	"fmt"
)

type Kademlia struct {
	routingTable *RoutingTable
  network *Network
}

func (kademlia *Kademlia) LookupContact(target *Contact) {

	const alpha := 3
	const k := 20

	var contacts []Contact
	var seenIDs []KademliaID

	done = make(chan []Contact)
	returned_contacts := make([]Contact, bucketSize)

  //pick out alpha (3) nodes from its closest non-empty k-bucket
  returned_contacts = append(returned_contacts, routingTable.FindClosestContacts(target.ID, alpha))

  // send parallel, async FIND_NODE to the alpha (3) nodes chosen
	findNode := 0	// amount of current findNode calls
	findNodeResponses := 0

  for _, contact := range returned_contacts {
		findNode++
		seenIDs = append(seenIDs, contact.ID)
    go network.SendFindContactMessage(contact, target, done)
  }

	// if a round of FIND_NODES fails to return a node any closer than the closest
	// already seen, the initiator resends the FIND_NODE to all k-closest nodes
	// it has not already queried

	// The lookup terminates when the initiator has queried and gotten responses
	// from the k closest nodes it has seen.

	for findNode > 0 && findNodeResponses > k{
		nodes := <-done
		findNodeResponses++
		findNode--

		for _, node := range nodes {
			// checks if node has already been seen, if so, ignore
			if !contains(seenIDs.String(), node.ID.String()) {
				routingTable.AddContact(node)
				returned_contacts = append(returned_contacts, node)
				seenIDs = append(seenIDs, node.ID)
			}
		}

		// resends FIND_NODE to nodes that it has learned from previous RPC
		for findNode < alpha && findNodeResponses > k{

			// sends new RPC to the closest contact in the priority queue that has not been seen
			contacts := routingTable.FindClosestContacts(target, k)

			// creates async FIND_NODE to the first unseen contact in the routingTable
			for _, contact := range contacts{
				if !contains(seenIDs.String(), contact.ID){
					findNode++
					go network.SendFindContactMessage(contact, target, done)
					break
				}
			}

		}

	}

}

// helper function for finding out if an element exists in a slice
// from stackoverflow https://stackoverflow.com/questions/10485743/contains-method-for-a-slice
func contains(slice []string, item string) bool {
    set := make(map[string]struct{}, len(slice))
    for _, s := range slice {
        set[s] = struct{}{}
    }

    _, ok := set[item]
    return ok
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
