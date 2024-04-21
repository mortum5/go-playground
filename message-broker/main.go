package main

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

// PUT
//   Check Request Queue
//     Send to request writer
//     Store in Topics
// Get
//   Check topicks
//   Timeout
//     Store in request queue
//

type RequestQueue struct {
	lock *sync.Mutex
	data map[string]*DoubleLinkedList
}

type DoubleLinkedList struct {
	head *Node
	tail *Node
}

type Node struct {
	prev *Node
	next *Node
	ch   chan string
}

func (rQueue *RequestQueue) Add(key string, ch chan string) *Node {
	rQueue.lock.Lock()
	defer rQueue.lock.Unlock()
	dll, ok := rQueue.data[key]
	if !ok {
		dll = &DoubleLinkedList{
			head: nil,
			tail: nil,
		}
		rQueue.data[key] = dll
	}
	node := &Node{nil, nil, ch}
	if dll.head == nil && dll.tail == nil {
		dll.head = node
		dll.tail = node
	} else {
		dll.tail.next = node
		node.prev = dll.tail
		dll.tail = node
	}
	return node

}

func (rQueue *RequestQueue) Get(key string) *Node {
	rQueue.lock.Lock()
	defer rQueue.lock.Unlock()
	dll, ok := rQueue.data[key]
	if !ok {
		return nil
	}

	return dll.head
}

func (rQueue *RequestQueue) Delete(key string, node *Node) {
	rQueue.lock.Lock()
	defer rQueue.lock.Unlock()
	dll, ok := rQueue.data[key]
	if !ok {
		return
	}

	if dll.head != node {
		node.prev.next = node.next
	} else {
		dll.head = node.next
	}

	if dll.tail != node {
		node.next.prev = node.prev
	} else {
		dll.tail = node.prev
	}

}

type Topics struct {
	lock *sync.Mutex
	data map[string][]string
}

func (t *Topics) Add(key, value string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if _, ok := t.data[key]; !ok {
		t.data[key] = []string{}
	}
	t.data[key] = append(t.data[key], value)
}

func (t *Topics) Get(key string) string {
	t.lock.Lock()
	defer t.lock.Unlock()
	_, ok := t.data[key]
	val := ""
	if len(t.data[key]) >= 1 {
		val = t.data[key][0]
		t.data[key] = t.data[key][1:]
	}

	if ok && len(t.data[key]) == 0 {
		delete(t.data, key)
	}

	return val
}

func main() {
	t := &Topics{
		data: make(map[string][]string),
	}

	rQueue := &RequestQueue{
		data: make(map[string]*DoubleLinkedList),
	}

	http.HandleFunc("/", queueHandle(t, rQueue))
	http.ListenAndServe(":8080", nil)
}

func queueHandle(t *Topics, rQueue *RequestQueue) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			queueName := req.URL.Path
			if queueName == "" {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			queueName = strings.TrimPrefix(queueName, "/")
			timeout := req.URL.Query().Get("timeout")
			value := t.Get(queueName)
			if value != "" {
				w.Write([]byte(value))
				return
			}
			if timeout == "" {
				http.NotFound(w, req)
				return
			}
			dur, err := time.ParseDuration(timeout + "s")
			if err != nil {
				http.Error(w, "", http.StatusBadRequest)
				return
			}
			ch := make(chan string)
			node := rQueue.Add(queueName, ch)
			defer rQueue.Delete(queueName, node)
			select {
			case msg := <-ch:
				w.Write([]byte(msg))
			case <-time.After(dur):
				http.NotFound(w, req)
			}
			return
		case "PUT":
			queueName := req.URL.Path
			if queueName == "" {
				http.Error(w, "", http.StatusBadRequest)
				return
			}
			queueName = strings.TrimPrefix(queueName, "/")
			value := req.URL.Query().Get("v")
			if value == "" {
				http.Error(w, "", http.StatusBadRequest)
				return
			}
			if node := rQueue.Get(queueName); node != nil {
				node.ch <- value
			} else {
				t.Add(queueName, value)
			}
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}

}
