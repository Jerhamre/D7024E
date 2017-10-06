package kademlia

import (
	"fmt"
)

type bucket struct {
	list []Contact
}

func newBucket() *bucket {
	var temp []Contact
	bucket := &bucket{temp}
	return bucket
}

func (bucket *bucket) AddContact(contact Contact) {
	/*var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
		}
	} else {
		bucket.list.MoveToFront(element)
	}*/
	fmt.Printf("AddContact %v\n", contact.String())

	var temp []Contact
	temp = append(temp, contact)

	index := contains(bucket.list, contact)
	fmt.Printf("index %v\n", index)


	if index != -1 {
		// Flytta längst fram
		bucket.list = append(bucket.list[:index], bucket.list[index+1:]...)
		bucket.list = append(temp, bucket.list...)
	} else {
		if len(bucket.list) < bucketSize {
			// Lägg till längst fram
			bucket.list = append(temp, bucket.list...)
		}
	}
}

func contains(s []Contact, e Contact) int {
  for i, a := range s {
		fmt.Println(e.String())
    if a.ID.String() == e.ID.String() {
      return i
    }
  }
  return -1
}

func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	/*var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts*/

	var contacts []Contact
	for _,c := range bucket.list {
		c.CalcDistance(target)
		contacts = append(contacts, c)
	}

	//fmt.Println("GetContactAndCalcDistance")
	//fmt.Println(contacts)
	return contacts
}

func (bucket *bucket) Len() int {
	return len(bucket.list)
}

func (bucket *bucket) Remove(contact Contact) {
	index := contains(bucket.list, contact)
	if index != -1 {
		bucket.list = append(bucket.list[:index], bucket.list[index+1:]...)
	}
}
