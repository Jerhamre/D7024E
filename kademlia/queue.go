package kademlia

import (
  "time"
)

type Queue struct {
  Waiting       chan Contact
  RoutingTable  *RoutingTable
}

func (queue *Queue) Run() {
  for {
    done := make(chan bool)
    go queue.Execute(done)
    <-done
    time.Sleep(time.Millisecond * 10)
  }
}

func (queue *Queue) Enqueue(c Contact) {
  queue.Waiting<-c
}

func (queue *Queue) Execute(done chan bool) {
  select {
    case c, ok := <- queue.Waiting:
      if ok {
        queue.RoutingTable.AddContact(c)
      }
    default:
    }
  done <- true
}