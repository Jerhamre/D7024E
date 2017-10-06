package kademlia

import (
  "time"
  "fmt"
)

type Queue struct {
  Waiting       chan Contact
  RoutingTable  *RoutingTable
  SleepTimer    time.Duration
}

func (queue *Queue) Run() {
  for {
    //fmt.Println(len(queue.Waiting))
    //if(len(queue.Waiting) > 0) {
      done := make(chan bool)
      go queue.Execute(done)
      <-done
    //}
    time.Sleep(time.Millisecond * queue.SleepTimer)
  }
}

func (queue *Queue) Enqueue(c Contact) {
  queue.Waiting<-c
}

func (queue *Queue) Execute(done chan bool) {
  select {
    case c, ok := <- queue.Waiting:
      if ok {
        fmt.Printf("Execute enqueue %v\n", c.String())
        if c.ID == nil {
          fmt.Println("Contact was nil. Did not add to routing table")
        } else {
          queue.RoutingTable.AddContact(c)
        }
      }
    default:
    }
  done <- true
}
