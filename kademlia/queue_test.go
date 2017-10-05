package kademlia

import (
	"testing"
  "time"
)

func TestEnqueue(t *testing.T) {

  kID := NewHashKademliaID("localhost:8000")
  me := NewContact(kID, "localhost:8000")
  me.CalcDistance(kID)
  rt := NewRoutingTable(me)
  queue := Queue{make(chan Contact, 100), rt, 1000}
  go queue.Run()

  c := NewContact(NewHashKademliaID("localhost:1111"), "localhost:1111")
  queue.Enqueue(c)

  if(len(queue.Waiting) != 1) {
    t.Fatalf("Expected %v in queue but got %v", 1, len(queue.Waiting))
  }
}

func TestExecute(t *testing.T) {

  kID := NewHashKademliaID("localhost:8000")
  me := NewContact(kID, "localhost:8000")
  me.CalcDistance(kID)
  rt := NewRoutingTable(me)
  queue := Queue{make(chan Contact, 100), rt, 10}
  go queue.Run()

  c := NewContact(NewHashKademliaID("localhost:1111"), "localhost:1111")
  queue.Enqueue(c)

  time.Sleep(time.Millisecond * 20)

  if(len(queue.Waiting) != 0) {
    t.Fatalf("Expected %v in queue but got %v", 0, len(queue.Waiting))
  }

  cs := queue.RoutingTable.FindClosestContacts(c.ID, 20)

  if(len(cs) != 1) {
    t.Fatalf("Expected %v contacts but got %v", 1, len(cs))
  }

  for _, c2 := range cs {
    if(c2.String() != c.String()) {
      t.Fatalf("Expected %v but got %v", c.String(), c2.String())
    }
  }
}
