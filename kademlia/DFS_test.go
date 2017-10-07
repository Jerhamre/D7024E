package kademlia

import (
	"testing"
  "reflect"
  "os"
  "time"
)
func TestNewDFS(t *testing.T) {

  os.RemoveAll("files")

  actualResult := NewDFS(nil, 60)
	var expectedResult DFS

  if(reflect.TypeOf(actualResult) != reflect.TypeOf(expectedResult)) {
    t.Fatalf("Expected %t but got %t", reflect.TypeOf(expectedResult), reflect.TypeOf(actualResult))
  }

  if(actualResult.PurgeTimer != 60) {
    t.Fatalf("Expected %v but got %v", 60, actualResult.PurgeTimer)
  }
}

func TestStore(t *testing.T) {

  kID := NewHashKademliaID("localhost:8000")
  me := NewContact(kID, "localhost:8000")
  me.CalcDistance(kID)
  rt := NewRoutingTable(me)
  dfs := NewDFS(rt, 10000)
  network := Network{"localhost", "8000"}
  queue := Queue{make(chan Contact), rt, 10}
  go queue.Run()
  k := Kademlia{rt, &network, &queue, &dfs}
  dfs.InitDFS(k)

  filename := "test1"
  data := []byte("file content")

  done := make(chan string)
  go dfs.Store(filename, data, done)
  back := <-done

	time.Sleep(time.Millisecond * 30)

  if(filename != back) {
    t.Fatalf("Expected filenames %v but got %v", filename, back)
  }

  if(!testEq(data, dfs.Files[filename].Data)) {
    t.Fatalf("Expected file data %v but got %v", data, dfs.Files[filename].Data)
  }
}

func TestCat(t *testing.T) {

  kID := NewHashKademliaID("localhost:8000")
  me := NewContact(kID, "localhost:8000")
  me.CalcDistance(kID)
  rt := NewRoutingTable(me)
  dfs := NewDFS(rt, 10000)
  network := Network{"localhost", "8000"}
  queue := Queue{make(chan Contact), rt, 10}
  go queue.Run()
  k := Kademlia{rt, &network, &queue, &dfs}
  dfs.InitDFS(k)

  filename := "test2"
  data := []byte("file content")

	store := make(chan string)
  go dfs.Store(filename, data, store)
  <-store

	time.Sleep(time.Millisecond * 30)

  done := make(chan []byte)
  go dfs.Cat(filename, done)
  back := <-done

  if(!testEq(data, back)) {
    t.Fatalf("Expected %v but got %v", data, back)
  }
}

func TestPin(t *testing.T) {

  kID := NewHashKademliaID("localhost:8000")
  me := NewContact(kID, "localhost:8000")
  me.CalcDistance(kID)
  rt := NewRoutingTable(me)
  dfs := NewDFS(rt, 40)
  network := Network{"localhost", "8000"}
  queue := Queue{make(chan Contact), rt, 10}
  go queue.Run()
  k := Kademlia{rt, &network, &queue, &dfs}
  dfs.InitDFS(k)

  filename := "test3"
  data := []byte("file content")

	store := make(chan string)
  go dfs.Store(filename, data, store)
  <-store

  time.Sleep(time.Millisecond * 20)

  done := make(chan bool)
  go dfs.Pin(filename, done)
  b := <-done

  if(b != true) {
    t.Fatalf("Expected true back from Pin %v but got %v", true, b)
  }

  time.Sleep(time.Millisecond * 50)

	cat := make(chan []byte)
  go dfs.Cat(filename, cat)
  back := <-cat

	if(!testEq(back, data)) {
    t.Fatalf("Expected file data %v but got %v", string(data), string(back))
  }
}

func TestUnpin(t *testing.T) {

  kID := NewHashKademliaID("localhost:8000")
  me := NewContact(kID, "localhost:8000")
  me.CalcDistance(kID)
  rt := NewRoutingTable(me)
  dfs := NewDFS(rt, 40)
  network := Network{"localhost", "8000"}
  queue := Queue{make(chan Contact), rt, 10}
  go queue.Run()
  k := Kademlia{rt, &network, &queue, &dfs}
  dfs.InitDFS(k)

  filename := "test4"
  data := []byte("file content")

	store := make(chan string)
  go dfs.Store(filename, data, store)
  <-store

  time.Sleep(time.Millisecond * 20)

  done := make(chan bool)
  go dfs.Pin(filename, done)
  b := <-done

  if(b != true) {
    t.Fatalf("Expected true back from Pin %v but got %v", true, b)
  }

	time.Sleep(time.Millisecond * 10)

	done2 := make(chan bool)
  go dfs.Unpin(filename, done2)
  b2 := <-done2

  if(b2 != true) {
    t.Fatalf("Expected true back from Unpin %v but got %v", true, b2)
  }

  time.Sleep(time.Millisecond * 55)

	expectedResult := []byte("CatFileDoesntExists")

	cat := make(chan []byte)
	go dfs.Cat(filename, cat)
	back := <-cat


	if(!testEq(back, expectedResult)) {
		t.Fatalf("Expected file data %v but got %v", string(expectedResult), back)
	}
}

func TestPurge(t *testing.T) {

  kID := NewHashKademliaID("localhost:8000")
  me := NewContact(kID, "localhost:8000")
  me.CalcDistance(kID)
  rt := NewRoutingTable(me)
  dfs := NewDFS(rt, 25)
  network := Network{"localhost", "8000"}
  queue := Queue{make(chan Contact), rt, 10}
  go queue.Run()
  k := Kademlia{rt, &network, &queue, &dfs}
  dfs.InitDFS(k)

  filename := "test5"
  data := []byte("file content")

	store := make(chan string)
  go dfs.Store(filename, data, store)
  <-store

  time.Sleep(time.Millisecond * 50)

	expectedResult := []byte("CatFileDoesntExists")

	cat := make(chan []byte)
	go dfs.Cat(filename, cat)
	back := <-cat


	if(!testEq(back, expectedResult)) {
		t.Fatalf("Expected file data %v but got %v", string(expectedResult), back)
	}
}

func TestStopPurge(t *testing.T) {

  kID := NewHashKademliaID("localhost:8000")
  me := NewContact(kID, "localhost:8000")
  me.CalcDistance(kID)
  rt := NewRoutingTable(me)
  dfs := NewDFS(rt, 40)
  network := Network{"localhost", "8000"}
  queue := Queue{make(chan Contact), rt, 10}
  go queue.Run()
  k := Kademlia{rt, &network, &queue, &dfs}
  dfs.InitDFS(k)

  filename := "test6"
  data := []byte("file content")

	store := make(chan string)
  go dfs.Store(filename, data, store)
  <-store

  time.Sleep(time.Millisecond * 20)

  go dfs.StopPurge(filename)

  time.Sleep(time.Millisecond * 50)

	cat := make(chan []byte)
	go dfs.Cat(filename, cat)
	back := <-cat

	if(!testEq(back, data)) {
		t.Fatalf("Expected file data %v but got %v", string(data), string(back))
	}
}
func TestGetFiles(t *testing.T) {

  kID := NewHashKademliaID("localhost:8000")
  me := NewContact(kID, "localhost:8000")
  me.CalcDistance(kID)
  rt := NewRoutingTable(me)
  dfs := NewDFS(rt, 5000)
  network := Network{"localhost", "8000"}
  queue := Queue{make(chan Contact), rt, 10}
  go queue.Run()
  k := Kademlia{rt, &network, &queue, &dfs}
  dfs.InitDFS(k)

  filename := "test7"
  data := []byte("file content")

  done := make(chan string)
  go dfs.Store(filename, data, done)
  <-done

  time.Sleep(time.Millisecond * 30)

  files := dfs.GetFiles()

  if(len(files) != 1) {
    t.Fatalf("Expected length %v but got %v", 1, len(files))
  }

  if(filename != files[0].Filename) {
    t.Fatalf("Expected filename %v but got %v", filename, files[0].Filename)
  }

  if(string(data) != files[0].Data) {
    t.Fatalf("Expected file content %v but got %v", string(data), files[0].Data)
  }
}

func testEq(a, b []byte) bool {
  if a == nil && b == nil {
    return true;
  }
  if a == nil || b == nil {
    return false;
  }
  if len(a) != len(b) {
    return false
  }
  for i := range a {
    if a[i] != b[i] {
      return false
    }
  }
  return true
}
