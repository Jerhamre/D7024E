package kademlia

import (
	"net/http"
	"html/template"
  "fmt"
  "encoding/json"
	"log"
)

type Page struct {
    Title string
    ID  string
    Address string
    ContactsCount  int
    Contacts []string
    ResultType string
    ResultFilename string
    ResultContent string
}

type Resonse struct {
  Filename string
  Content string
}

func HTTPListen(port string, kademlia *Kademlia) {
  http.HandleFunc("/", indexHandler(kademlia))
  http.HandleFunc("/store", httpStore(kademlia))
  http.HandleFunc("/cat", httpCat(kademlia))
  http.HandleFunc("/pin", httpPin(kademlia))
  http.HandleFunc("/unpin", httpUnpin(kademlia))

	fmt.Println("Listening HTTP server on:\t" + kademlia.Network.IP + ":" + port)

  http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("kademlia/templates/resources"))))
	err := http.ListenAndServe(":"+port,nil)
	if err != nil {
  	log.Fatal("ListenAndServe: ", err)
	}
}

func getPage(kademlia *Kademlia) *Page {
  s := []string{}
  i := 0


  for _, c := range kademlia.RoutingTable.FindClosestContacts(kademlia.RoutingTable.me.ID, 160) {
    s = append(s, c.Address + " " + c.ID.String())
    i++
  }

  p := &Page{
    Title: "Kademlia",
    ID: kademlia.RoutingTable.me.ID.String(),
    Address: kademlia.RoutingTable.me.Address,
    ContactsCount: i,
    Contacts: s,
  }
  return p
}

func indexHandler(kademlia *Kademlia) func ( w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request) {
    p := getPage(kademlia)
    t,err := template.ParseFiles("kademlia/templates/index.html")
		if err != nil { panic(err) }
    err = t.Execute(w, p)
		if err != nil { panic(err) }
	}
}
func httpStore(kademlia *Kademlia) func ( w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request) {
    fmt.Println("store")
    if r.Method == "POST" {
      decoder := json.NewDecoder(r.Body)

      var jmsg Resonse
      err := decoder.Decode(&jmsg)
      if err != nil {
        fmt.Println(err)
      }
      filename := jmsg.Filename
      content := jmsg.Content

			done := make(chan string)
			go kademlia.Store(filename, []byte(content), done)

      fmt.Fprint(w, <-done);
    } else {
      fmt.Fprint(w, "0");
    }
	}
}
func httpCat(kademlia *Kademlia) func ( w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request) {
    fmt.Println("cat")
    if r.Method == "GET" {
      decoder := json.NewDecoder(r.Body)

      var jmsg Resonse
      err := decoder.Decode(&jmsg)
      if err != nil {
        fmt.Println(err)
      }
      filename := jmsg.Filename

			done := make(chan []byte)
			go kademlia.LookupData(filename, done)

      fmt.Fprint(w, string(<-done));
    } else {
      fmt.Fprint(w, "0");
    }
	}
}
func httpPin(kademlia *Kademlia) func ( w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request) {
    fmt.Println("pin")
    if r.Method == "POST" {
      decoder := json.NewDecoder(r.Body)

      var jmsg Resonse
      err := decoder.Decode(&jmsg)
      if err != nil {
        fmt.Println(err)
      }
      filename := jmsg.Filename

      fmt.Println(filename)

			done := make(chan bool)
			go kademlia.Pin(filename, done)

      fmt.Fprint(w, <-done);
    } else {
      fmt.Fprint(w, "0");
    }
	}
}
func httpUnpin(kademlia *Kademlia) func ( w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request) {
    fmt.Println("unpin")
    if r.Method == "POST" {
      decoder := json.NewDecoder(r.Body)

      var jmsg Resonse
      err := decoder.Decode(&jmsg)
      if err != nil {
        fmt.Println(err)
      }
      filename := jmsg.Filename

      fmt.Println(filename)

			done := make(chan bool)
			go kademlia.Unpin(filename, done)

      fmt.Fprint(w, <-done);
    } else {
      fmt.Fprint(w, "0");
    }
	}
}
