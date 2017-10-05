package kademlia

import (
	"testing"
  "net/http"
  "io/ioutil"
  "strings"
	"bytes"
)

func TestHTTPIndex(t *testing.T) {
  kID := NewHashKademliaID("localhost:9000")
  me := NewContact(kID, "localhost:9000")
  me.CalcDistance(kID)
  rt := NewRoutingTable(me)
  dfs := NewDFS(rt, 10000)
  network := Network{"localhost", "9000"}
  queue := Queue{make(chan Contact), rt, 10}
  go queue.Run()
  k := Kademlia{rt, &network, &queue, &dfs}
  go HTTPListen("9000", &k, "templates")

  resp, err := http.Get("http://localhost:9000")
  err = err
  resp = resp
  bytes, _ := ioutil.ReadAll(resp.Body)
  resp.Body.Close()

  if(!strings.Contains(string(bytes), "6c6f63616c686f73743a39303030202020202020")) {
    t.Fatalf("Expected %v  in HTML response but got %v", "6c6f63616c686f73743a39303030202020202020", string(bytes))
  }
}

func TestHTTPStore(t *testing.T) {

		url := "http://localhost:9000/store"

		filename := "test"
		data := []byte("file content")

		var jsonStr = []byte("{\"Filename\":\""+filename+"\", \"Content\": \""+string(data)+"\"}")
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	  req.Header.Set("Content-Type", "application/json")

	  client := &http.Client{}
	  resp, err := client.Do(req)
	  if err != nil {
	    panic(err)
	  }
	  defer resp.Body.Close()


		body, _ := ioutil.ReadAll(resp.Body)

	  if(string(body) != filename) {
	    t.Fatalf("Expected %v but got %v", filename, string(body))
		}
}

func TestHTTPCat(t *testing.T) {

	url := "http://localhost:9000/store"

	filename := "test"
	data := []byte("file content")

	var jsonStr = []byte("{\"Filename\":\""+filename+"\", \"Content\": \""+string(data)+"\"}")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()


	body, _ := ioutil.ReadAll(resp.Body)

	if(string(body) != filename) {
		t.Fatalf("Expected Store %v but got %v", filename, string(body))
	}

	url2 := "http://localhost:9000/cat?Filename="+filename

	req2, err2 := http.NewRequest("GET", url2, nil)

	client2 := &http.Client{}
	resp2, err2 := client2.Do(req2)
	if err2 != nil {
		panic(err2)
	}
	defer resp2.Body.Close()


	body2, _ := ioutil.ReadAll(resp2.Body)

	if(string(body2) != string(data)) {
		t.Fatalf("Expected Cat %v but got %v", string(data), string(body2))
	}
}

func TestHTTPPin(t *testing.T) {

	url := "http://localhost:9000/store"

	filename := "test"
	data := []byte("file content")

	var jsonStr = []byte("{\"Filename\":\""+filename+"\", \"Content\": \""+string(data)+"\"}")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()


	body, _ := ioutil.ReadAll(resp.Body)

	if(string(body) != filename) {
		t.Fatalf("Expected Store %v but got %v", filename, string(body))
	}

	url2 := "http://localhost:9000/pin"

	var jsonStr2 = []byte("{\"Filename\":\""+filename+"\"}")
	req2, err2 := http.NewRequest("POST", url2, bytes.NewBuffer(jsonStr2))
	req2.Header.Set("Content-Type", "application/json")

	client2 := &http.Client{}
	resp2, err2 := client2.Do(req2)
	if err2 != nil {
		panic(err2)
	}
	defer resp2.Body.Close()


	body2, _ := ioutil.ReadAll(resp2.Body)

	if(string(body2) != "true") {
		t.Fatalf("Expected Pin %v but got %v", "true", string(body2))
	}
}

func TestHTTPUnpin(t *testing.T) {
	url := "http://localhost:9000/store"

	filename := "test"
	data := []byte("file content")

	var jsonStr = []byte("{\"Filename\":\""+filename+"\", \"Content\": \""+string(data)+"\"}")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()


	body, _ := ioutil.ReadAll(resp.Body)

	if(string(body) != filename) {
		t.Fatalf("Expected Store %v but got %v", filename, string(body))
	}

	url2 := "http://localhost:9000/unpin"

	var jsonStr2 = []byte("{\"Filename\":\""+filename+"\"}")
	req2, err2 := http.NewRequest("POST", url2, bytes.NewBuffer(jsonStr2))
	req2.Header.Set("Content-Type", "application/json")

	client2 := &http.Client{}
	resp2, err2 := client2.Do(req2)
	if err2 != nil {
		panic(err2)
	}
	defer resp2.Body.Close()


	body2, _ := ioutil.ReadAll(resp2.Body)

	if(string(body2) != "true") {
		t.Fatalf("Expected Unpin %v but got %v", "true", string(body2))
	}
}
