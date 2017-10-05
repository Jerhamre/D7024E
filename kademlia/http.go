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
    Contacts []string
    Files []File
    ResultType string
    ResultFilename string
    ResultContent string
}

type Response struct {
  Filename string
  Content string
}

func HTTPListen(port string, kademlia *Kademlia, templatePath string) {
  http.HandleFunc("/", indexHandler(kademlia, templatePath))
  http.HandleFunc("/store", httpStore(kademlia))
  http.HandleFunc("/cat", httpCat(kademlia))
  http.HandleFunc("/pin", httpPin(kademlia))
  http.HandleFunc("/unpin", httpUnpin(kademlia))

  http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir(templatePath+"/resources"))))

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
  	log.Fatal("ListenAndServe: ", err)
	}
}

func getPage(kademlia *Kademlia) *Page {
  contacts := []string{}
  for _, c := range kademlia.RoutingTable.FindClosestContacts(kademlia.RoutingTable.me.ID, 200) {
    contacts = append(contacts, c.Address + " " + c.ID.String())
  }


  p := &Page{
    Title: "Kademlia",
    ID: kademlia.RoutingTable.me.ID.String(),
    Address: kademlia.RoutingTable.me.Address,
    Contacts: contacts,
    Files: kademlia.DFS.GetFiles(),
  }
  return p
}

func indexHandler(kademlia *Kademlia, templatePath string) func ( w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request) {
    p := getPage(kademlia)
    t,err := template.ParseFiles(templatePath+"/index.html")
		if err != nil { panic(err) }
    err = t.Execute(w, p)
		if err != nil { panic(err) }
	}
}
func httpStore(kademlia *Kademlia) func ( w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
      decoder := json.NewDecoder(r.Body)

      var jmsg Response
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
    if r.Method == "GET" {
      queryValues := r.URL.Query()

      filename := queryValues.Get("Filename")

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
    if r.Method == "POST" {
      decoder := json.NewDecoder(r.Body)

      var jmsg Response
      err := decoder.Decode(&jmsg)
      if err != nil {
        fmt.Println(err)
      }
      filename := jmsg.Filename

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
    if r.Method == "POST" {
      decoder := json.NewDecoder(r.Body)

      var jmsg Response
      err := decoder.Decode(&jmsg)
      if err != nil {
        fmt.Println(err)
      }
      filename := jmsg.Filename

			done := make(chan bool)
			go kademlia.Unpin(filename, done)

      fmt.Fprint(w, <-done);
    } else {
      fmt.Fprint(w, "0");
    }
	}
}
