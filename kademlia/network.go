package kademlia

import (
  "fmt"
  "net"
  "encoding/json"
  "math/rand"
  "time"
  "os"
  "bufio"
  "strconv"
  "strings"
)

type Network struct {
  IP  string
  Port string
}

type Message struct {
  RCP_ID int
  KademliaID *KademliaID
  Address  string
  MessageType string
  Data string
}


/* A Simple function to verify error */
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}
func Listen(kademlia *Kademlia) {
  port, err := strconv.Atoi(kademlia.Network.Port)

  in := make([]byte, 8096)
  addr := net.UDPAddr{
    Port: port,
    IP: net.ParseIP("0.0.0.0"),
  }

  conn, err := net.ListenUDP("udp", &addr)
  if err != nil {
    fmt.Printf("Some error %v\n", err)
    return
  }

  //fmt.Println("Listening Kademlia on:\t\t" + kademlia.Network.IP + ":" + kademlia.Network.Port)

  for {
    n,remoteaddr,err := conn.ReadFromUDP(in)
    if err !=  nil {
      fmt.Printf("Some error  %v", err)
      continue
    }
    go handleRequest(conn, in[:n], remoteaddr, kademlia)
  }
}

func (network *Network) SendPingMessage(me *Contact, contact *Contact, done chan *KademliaID, iteration ... int) {
  p :=  make([]byte, 8096)

  var errorRes *KademliaID

  var iter int
  if iteration == nil {
    iter = 0
  } else {
    iter = iteration[0]
  }
  if(iter >= 5) {
    fmt.Println("Too many tries in SendPingMessage...")
    done<-errorRes
    return
  }

  conn, err := net.Dial("udp", contact.Address)
  if err != nil {
    fmt.Printf("Some error SPM 1 %v", err)
    done<-errorRes
    return
  }

  RCP_ID := getRandomRCPID()
  out := packMessage(me, RCP_ID, "SendPingMessage", "")

  fmt.Fprintf(conn, out)

  stopTimeout := make(chan bool)
  go func(network *Network, me *Contact, contact *Contact, done chan *KademliaID) {
    select {
      case <-stopTimeout:
        break
      case <-time.After(2 * time.Second):
        conn.Close()
        fmt.Println("SendPingMessage timed out", iter, ", recalling again...")
        network.SendPingMessage(me, contact, done, iter+1)
        return
    }
  }(network, me, contact, done)

  n, err := bufio.NewReader(conn).Read(p)

  stopTimeout <- true

  if err == nil {
    in := unpackMessage(string(p[:n]))
    done <-NewKademliaID(in.Data)
  } else {
    fmt.Printf("Some error SPM 2 %v\n", err)
    done<-errorRes
    return
  }
  conn.Close()
}

func (network *Network) SendFindContactMessage(me *Contact, contact *Contact, target *Contact, done chan []Contact, iteration ... int) {
  p :=  make([]byte, 8096)

  var errorRes []Contact

  var iter int
  if iteration == nil {
    iter = 0
  } else {
    iter = iteration[0]
  }
  if(iter >= 5) {
    fmt.Println("Too many tries in SendFindContactMessage...")
    done<-errorRes
    return
  }

  conn, err := net.Dial("udp", contact.Address)
  if err != nil {
    fmt.Printf("Some error SFCM 1 %v", err)
    done<-errorRes
    return
  }

  RCP_ID := getRandomRCPID()
  targetData, err := json.Marshal(target)
  if err != nil {
    fmt.Println(err)
  }
  out := packMessage(me, RCP_ID, "SendFindContactMessage", string(targetData))

  fmt.Fprintf(conn, out)

  stopTimeout := make(chan bool)
  go func(network *Network, me *Contact, contact *Contact, target *Contact, done chan []Contact) {
    select {
      case <-stopTimeout:
        break
      case <-time.After(2 * time.Second):
        conn.Close()
        fmt.Println("SendFindContactMessage timed out", iter, ", recalling again...")
        network.SendFindContactMessage(me, contact, target, done, iter+1)
        return
    }
  }(network, me, contact, target, done)

  n, err := bufio.NewReader(conn).Read(p)

  stopTimeout <- true

  if err == nil {
    in := unpackMessage(string(p[:n]))
    var contacts []Contact
    err = json.Unmarshal([]byte(in.Data), &contacts)

    if err != nil {
      fmt.Println(err)
    }
  done<-contacts
  } else {
    fmt.Printf("Some error SFCM 2%v\n", err)
    done<-errorRes
    return
  }
  conn.Close()
}

func (network *Network) SendFindDataMessage(me *Contact, contact *Contact, filename string, done chan []byte, iteration ... int) {
  p :=  make([]byte, 8096)

  var errorRes []byte

  var iter int
  if iteration == nil {
    iter = 0
  } else {
    iter = iteration[0]
  }
  if(iter >= 5) {
    fmt.Println("Too many tries in SendFindDataMessage...")
    done<-errorRes
    return
  }

  conn, err := net.Dial("udp", contact.Address)
  if err != nil {
    fmt.Printf("Some error %v", err)
    done <- errorRes
    return
  }

  RCP_ID := getRandomRCPID()
  out := packMessage(me, RCP_ID, "SendFindDataMessage", filename)

  fmt.Fprintf(conn, out)

  stopTimeout := make(chan bool)
  go func(network *Network, me *Contact, contact *Contact, filename string, done chan []byte) {
    select {
      case <-stopTimeout:
        break
      case <-time.After(2 * time.Second):
        conn.Close()
        fmt.Println("SendFindDataMessage timed out", iter, ", recalling again...")
        network.SendFindDataMessage(me, contact, filename, done, iter+1)
        return
    }
  }(network, me, contact, filename, done)

  n, err := bufio.NewReader(conn).Read(p)

  stopTimeout <- true

  if err == nil {
    in := unpackMessage(string(p[:n]))
    done<-[]byte(in.Data)
  } else {
    fmt.Printf("Some error %v\n", err)
    done <- errorRes
  }
  conn.Close()
}

func (network *Network) SendStoreMessage(me *Contact, contact *Contact, filename string, data []byte, done chan string, iteration ... int) {
  p :=  make([]byte, 8096)

  var errorRes string

  var iter int
  if iteration == nil {
    iter = 0
  } else {
    iter = iteration[0]
  }
  if(iter >= 5) {
    fmt.Println("Too many tries in SendStoreMessage...")
    done<-errorRes
    return
  }

  conn, err := net.Dial("udp", contact.Address)
  if err != nil {
    fmt.Printf("Some error SendStoreMessage %v", err)
    done <- "fail"
    return
  }

  RCP_ID := getRandomRCPID()
  out := packMessage(me, RCP_ID, "SendStoreMessage", strings.Join([]string{filename, string(data)}, ","))

  fmt.Fprintf(conn, out)

  stopTimeout := make(chan bool)
  go func(network *Network, me *Contact, contact *Contact, filename string, data []byte, done chan string) {
    select {
      case <-stopTimeout:
        break
      case <-time.After(2 * time.Second):
        conn.Close()
        fmt.Println("SendStoreMessage timed out", iter, ", recalling again...")
        network.SendStoreMessage(me, contact, filename, data, done, iter+1)
        return
    }
  }(network, me, contact, filename, data, done)

  n, err := bufio.NewReader(conn).Read(p)

  stopTimeout <- true

  if err == nil {
    in := unpackMessage(string(p[:n]))
    done <- in.Data
  } else {
    fmt.Printf("Some error Unpacking %v\n", err)
    done <- "fail"
  }
  conn.Close()
}

func (network *Network) SendPinMessage(me *Contact, contact *Contact, filename string, done chan bool, iteration ... int) {
  p :=  make([]byte, 8096)

  var errorRes bool

  var iter int
  if iteration == nil {
    iter = 0
  } else {
    iter = iteration[0]
  }
  if(iter >= 5) {
    fmt.Println("Too many tries in SendPinMessage...")
    done<-errorRes
    return
  }

  conn, err := net.Dial("udp", contact.Address)
  if err != nil {
    fmt.Printf("Some error %v", err)
    done <- errorRes
    return
  }

  RCP_ID := getRandomRCPID()
  out := packMessage(me, RCP_ID, "SendPinMessage", filename)

  fmt.Fprintf(conn, out)

  stopTimeout := make(chan bool)
  go func(network *Network, me *Contact, contact *Contact, filename string, done chan bool) {
    select {
      case <-stopTimeout:
        break
      case <-time.After(2 * time.Second):
        conn.Close()
        fmt.Println("SendPinMessage timed out", iter, ", recalling again...")
        network.SendPinMessage(me, contact, filename, done, iter+1)
        return
    }
  }(network, me, contact, filename, done)

  n, err := bufio.NewReader(conn).Read(p)

  stopTimeout <- true

  if err == nil {
    in := unpackMessage(string(p[:n]))
    b, err := strconv.ParseBool(in.Data)
    if err != nil {
      panic("Not a bool in SendPinMessage?")
    }
    done<-b
  } else {
    fmt.Printf("Some error %v\n", err)
    done <- errorRes
  }
  conn.Close()
}

func (network *Network) SendUnpinMessage(me *Contact, contact *Contact, filename string, done chan bool, iteration ... int) {
  p :=  make([]byte, 8096)

  var errorRes bool

  var iter int
  if iteration == nil {
    iter = 0
  } else {
    iter = iteration[0]
  }
  if(iter >= 5) {
    fmt.Println("Too many tries in SendUnpinMessage...")
    done<-errorRes
    return
  }

  conn, err := net.Dial("udp", contact.Address)
  if err != nil {
    fmt.Printf("Some error %v", err)
    done <- errorRes
    return
  }

  RCP_ID := getRandomRCPID()
  out := packMessage(me, RCP_ID, "SendUnpinMessage", filename)

  fmt.Fprintf(conn, out)

  stopTimeout := make(chan bool)
  go func(network *Network, me *Contact, contact *Contact, filename string, done chan bool) {
    select {
      case <-stopTimeout:
        break
      case <-time.After(2 * time.Second):
        conn.Close()
        fmt.Println("SendUnpinMessage timed out", iter, ", recalling again...")
        network.SendUnpinMessage(me, contact, filename, done, iter+1)
        return
    }
  }(network, me, contact, filename, done)

  n, err := bufio.NewReader(conn).Read(p)

  stopTimeout <- true

  if err == nil {
    in := unpackMessage(string(p[:n]))
    b, err := strconv.ParseBool(in.Data)
    if err != nil {
      panic("Not a bool in SendUnpinMessage?")
    }
    done<-b
  } else {
    fmt.Printf("Some error %v\n", err)
    done <- errorRes
  }
  conn.Close()
}

func getRandomRCPID() int {
    rand.Seed(time.Now().UnixNano())
    return rand.Int()
}

func packMessage(me *Contact, RCP_ID int, messageType string, data string) string {
  message := new(Message)
  message.RCP_ID = RCP_ID
  message.KademliaID = me.ID
  message.Address = me.Address
  message.MessageType = messageType
  message.Data = data
  out, err := json.Marshal(message)
  if err != nil {
    fmt.Println(err)
  }
  return string(out)
}

func unpackMessage(message string) Message {
  in := []byte(message)
  var raw Message //map[string]interface{}
  err := json.Unmarshal(in, &raw)
  if err != nil {
    fmt.Printf("UnpackMessage: %v\n", err)
    fmt.Println(in)
    fmt.Println(string(in))
  }
  return raw
}

func handleRequest(conn *net.UDPConn, buf []byte, remoteaddr *net.UDPAddr, kademlia *Kademlia) {
  in := unpackMessage(string(buf))

  sender := NewContact(in.KademliaID, in.Address)
  data := ""

  switch in.MessageType {
    case "SendPingMessage":
      data = kademlia.RoutingTable.me.ID.String()
    case "SendFindContactMessage":
      var target Contact //map[string]interface{}
      err := json.Unmarshal([]byte(in.Data), &target)
      if err != nil {
        fmt.Println(err)
      }

      kademlia.Queue.Enqueue(sender)

      contacts := kademlia.RoutingTable.FindClosestContacts(target.ID, 20)

      out, err := json.Marshal(contacts)
      if err != nil {
        fmt.Println(err)
      }
      data = string(out)

    case "SendFindDataMessage":
      hash := in.Data

      done := make(chan []byte)
      go kademlia.DFS.Cat(hash, done)

      data = string(<-done)

    case "SendStoreMessage":
      s := strings.Split(in.Data, ",")
      filename := s[0]
      content := s[1]

      done := make(chan string)
      go kademlia.DFS.Store(filename, []byte(content), done)

      data = <-done
    case "SendPinMessage":
      hash := in.Data

      done := make(chan bool)
      go kademlia.DFS.Pin(hash, done)

      data = strconv.FormatBool(<-done)
    case "SendUnpinMessage":
      hash := in.Data

      done := make(chan bool)
      go kademlia.DFS.Unpin(hash, done)

      data = strconv.FormatBool(<-done)
    default:
      fmt.Println("Not a valid message type")
      fmt.Println(in)
      return
  }

  if data == "" {
    panic("Not valid return data")
    conn.Close()
    return
  }

  out := packMessage(&kademlia.RoutingTable.me, in.RCP_ID, in.MessageType, data)

  // Send a response back to person contacting us.
  _,err := conn.WriteToUDP([]byte(out), remoteaddr)
  if err != nil {
      fmt.Printf("Couldn't send response %v", err)
  }
}
