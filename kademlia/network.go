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

    in := make([]byte, 2048)
    addr := net.UDPAddr{
        Port: port,
        IP: net.ParseIP("0.0.0.0"),
    }
    fmt.Println(addr)
    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Printf("Some error %v\n", err)
        return
    }

    fmt.Println("Listening Kademlia on:\t\t" + kademlia.Network.IP + ":" + kademlia.Network.Port)

    for {
        n,remoteaddr,err := conn.ReadFromUDP(in)
        if err !=  nil {
            fmt.Printf("Some error  %v", err)
            continue
        }
        fmt.Println("handleRequest")
        go handleRequest(conn, in[:n], remoteaddr, kademlia)
    }
}


func (network *Network) SendPingMessage(me *Contact, contact *Contact, done chan *KademliaID) {
  fmt.Println("SendPingMessage")
  p :=  make([]byte, 2048)

  var errorRes *KademliaID

  conn, err := net.Dial("udp", contact.Address)
  if err != nil {
    fmt.Printf("Some error SPM 1 %v", err)
    done<-errorRes
    return
  }

  fmt.Println(contact.Address)

  RCP_ID := getRandomRCPID()
  out := packMessage(me, RCP_ID, "SendPingMessage", "")

  fmt.Fprintf(conn, out)

   n, err := bufio.NewReader(conn).Read(p)

  if err == nil {
    in := unpackMessage(string(p[:n]))
    done <-NewKademliaID(in.Data)
  } else {
    fmt.Printf("Some error SPM 2 %v\n", err)
    done<-errorRes
  }
  conn.Close()
}

func (network *Network) SendFindContactMessage(me *Contact, contact *Contact, target *Contact, done chan []Contact) {

  fmt.Println("SendFindContactMessage")
  var errorRes []Contact
  p :=  make([]byte, 2048)
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

  n, err := bufio.NewReader(conn).Read(p)
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

func (network *Network) SendFindDataMessage(me *Contact, contact *Contact, filename string, done chan []byte) {
  fmt.Println("SendFindDataMessage")
  p :=  make([]byte, 2048)
  conn, err := net.Dial("udp", contact.Address)
  if err != nil {
    fmt.Printf("Some error %v", err)
    done <- []byte("fail")
    return
  }

  RCP_ID := getRandomRCPID()
  out := packMessage(me, RCP_ID, "SendFindDataMessage", filename)

  fmt.Fprintf(conn, out)

  n, err := bufio.NewReader(conn).Read(p)
  if err == nil {
    in := unpackMessage(string(p[:n]))
    done<-[]byte(in.Data)
  } else {
    fmt.Printf("Some error %v\n", err)
    done <- []byte("fail")
  }
  conn.Close()
}

func (network *Network) SendStoreMessage(me *Contact, contact *Contact, filename string, data []byte, done chan string) {
  fmt.Println("SendStoreMessage")
  p :=  make([]byte, 2048)
  conn, err := net.Dial("udp", contact.Address)
  if err != nil {
    fmt.Printf("Some error %v", err)
    done <- "fail"
    return
  }

  RCP_ID := getRandomRCPID()
  out := packMessage(me, RCP_ID, "SendStoreMessage", strings.Join([]string{filename, string(data)}, ","))

  fmt.Fprintf(conn, out)

  n, err := bufio.NewReader(conn).Read(p)

  if err == nil {
    in := unpackMessage(string(p[:n]))
    done <- in.Data
  } else {
    fmt.Printf("Some error %v\n", err)
    done <- "fail"
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
    fmt.Println(err)
  }
  return raw
}

func handleRequest(conn *net.UDPConn, buf []byte, remoteaddr *net.UDPAddr, kademlia *Kademlia) {
  in := unpackMessage(string(buf))

  fmt.Println(string(buf))

  sender := NewContact(in.KademliaID, in.Address)
  data := ""

  switch in.MessageType {
  case "SendPingMessage":
      fmt.Println("handle SendPingMessage")
      data = kademlia.RoutingTable.me.ID.String()
    case "SendFindContactMessage":
      fmt.Println("handle SendFindContactMessage")
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
      fmt.Println("handle SendFindDataMessage")
      hash := in.Data

      done := make(chan []byte)
      go kademlia.DFS.Cat(hash, done)

      data = string(<-done)

    case "SendStoreMessage":
      fmt.Println("handle SendStoreMessage")
      s := strings.Split(in.Data, ",")
      filename := s[0]
      content := s[1]

      done := make(chan string)
      go kademlia.DFS.Store(filename, []byte(content), done)

      data = <-done

  default:
      panic("Not a valid message type")
      fmt.Println(in)
      conn.Close()
      return
  }

  if data == "" {
    panic("Not valid return data")
    conn.Close()
    return
  }

  out := packMessage(&kademlia.RoutingTable.me, in.RCP_ID, in.MessageType, data)

  // Send a response back to person contacting us.
  fmt.Println(remoteaddr)

  _,err := conn.WriteToUDP([]byte(out), remoteaddr)
  if err != nil {
      fmt.Printf("Couldn't send response %v", err)
  }
}
