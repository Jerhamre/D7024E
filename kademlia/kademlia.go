package kademlia

import (
	"fmt"
)

type Kademlia struct {
	RoutingTable  *RoutingTable
	Network				*Network
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	fmt.Println("LookupContact")
}

func (kademlia *Kademlia) LookupData(hash string) {
	fmt.Println("LookupData")
}

func (kademlia *Kademlia) Store(data []byte) {
	fmt.Println("Store")
}
