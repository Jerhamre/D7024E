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
	var temp []Contact
	temp = append(temp, contact)

	index := contains(bucket.list, contact)
	if index != -1 {
		bucket.list = append(bucket.list[:index], bucket.list[index+1:]...)
		bucket.list = append(temp, bucket.list...)
	} else {
		if len(bucket.list) < bucketSize {
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
	var contacts []Contact
	for _,c := range bucket.list {
		c.CalcDistance(target)
		contacts = append(contacts, c)
	}
	
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
