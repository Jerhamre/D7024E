package kademlia

import (
  "fmt"
  "os"
  "net"
  "bufio"
  "encoding/json"
  "math/rand"
  "time"
)

type Network struct {
  IP  string
  Port string
}

type Message struct {
  RCP_ID int
  IP  string
  Port string
  MessageType string
  Data string
}

func Listen(kademlia *Kademlia) {
  // Listen for incoming connections.
  //kademlia.Network.IP+
    l, err := net.Listen("tcp", ":"+kademlia.Network.Port)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()

    fmt.Println("Listening on " + kademlia.Network.IP + ":" + kademlia.Network.Port)

    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go handleRequest(conn, kademlia)
    }
}

func (network *Network) SendPingMessage(contact *Contact) {

  conn, err := net.Dial("tcp", contact.Address)
  if err != nil {
    println("Dial failed:", err.Error())
  } else {
    RCP_ID := getRandomID()
    out := packMessage(network, RCP_ID, "SendPingMessage", "")

    fmt.Println("Client sent:", out)
    // send to socket
    fmt.Fprintf(conn, string(out)+"\n")
    // listen for reply
    message, _ := bufio.NewReader(conn).ReadString('\n')
    in := unpackMessage(message)
    fmt.Println("Client received:", in)
  }
}

func (network *Network) SendFindContactMessage(contact *Contact, target *Contact, done chan []Contact) {
  conn, err := net.Dial("tcp", contact.Address)
  if err != nil {
    println("Dial failed:", err.Error())
  } else {
    RCP_ID := getRandomID()

    targetData, err := json.Marshal(target)
    if err != nil {
      fmt.Println(err)
    }
    out := packMessage(network, RCP_ID, "SendFindContactMessage", string(targetData))

    // send to socket
    fmt.Fprintf(conn, string(out)+"\n")
    // listen for reply
    message, _ := bufio.NewReader(conn).ReadString('\n')
    in := unpackMessage(message)
    fmt.Println("Client received:", in)


    var contacts []Contact
    err = json.Unmarshal([]byte(in.Data), &contacts)
    if err != nil {
      fmt.Println(err)
    }

    //fmt.Println(contacts[0].String())

    done <- contacts
  }
}

func (network *Network) SendFindDataMessage(hash string) {
  fmt.Println("SendFindDataMessage")
}

func (network *Network) SendStoreMessage(data []byte) {
  fmt.Println("SendStoreMessage")
}

func getRandomID() int {
    rand.Seed(time.Now().UnixNano())
    return rand.Int()
}

func packMessage(network *Network, RCP_ID int, messageType string, data string) string {
  message := new(Message)
  message.RCP_ID = RCP_ID
  message.IP = network.IP
  message.Port = network.Port
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

func handleRequest(conn net.Conn, kademlia *Kademlia) {

  // Make a buffer to hold incoming data.
  buf_ := make([]byte, 1024)
  // Get the length of the incoming data.
  len, err := conn.Read(buf_)
  if err != nil {
    fmt.Println("Error reading:", err.Error())
  }
  // Read the incoming connection into the buffer.
  buf := buf_[:len]

  message := string(buf)
  in := unpackMessage(message)
  //fmt.Println("Server received:", in)

  data := ""

  switch in.MessageType {
    case "SendPingMessage":
      // Return time of contant that is online
      data = time.Now().String()

    case "SendFindContactMessage":

      var target Contact //map[string]interface{}
      err := json.Unmarshal([]byte(in.Data), &target)
      if err != nil {
        fmt.Println(err)
      }

      contacts := kademlia.RoutingTable.FindClosestContacts(target.ID, 20)

      out, err := json.Marshal(contacts)
      if err != nil {
        fmt.Println(err)
      }
      data = string(out)

    case "SendFindDataMessage":
      // Do something

    case "SendStoreMessage":
      // Do something

  default:
      panic("Not a valid message type")
      conn.Close()
      return
  }

  if data == "" {
    panic("Not valid return data")
    conn.Close()
    return
  }

  out := packMessage(kademlia.Network, in.RCP_ID, in.MessageType, data)

  // Send a response back to person contacting us.
  //fmt.Println("Server sent:", out)
  conn.Write([]byte(out))
  // Close the connection when you're done with it.
  conn.Close()
}
