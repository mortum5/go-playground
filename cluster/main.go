package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/memberlist"
)

var (
	counter GCounter
)

type GCounter struct {
	Ident   string
	Counter map[string]int
}

type delegate struct{}

func (d *delegate) NodeMeta(limit int) []byte {
	return []byte("")
}

func (d *delegate) NotifyMsg(_ []byte) {

}

func (d *delegate) GetBroadcasts(overhead int, limit int) [][]byte {
	return nil
}

func (d *delegate) LocalState(join bool) []byte {
	b, err := json.Marshal(counter)
	slog.Info("local state request", "node", counter.Ident, "marshal", string(b))
	if err != nil {
		panic(err)
	}

	return b
}

func (d *delegate) MergeRemoteState(buf []byte, join bool) {
	if len(buf) == 0 {
		return
	}

	var extCount GCounter
	json.Unmarshal(buf, &extCount)
	slog.Info("merge remote state request", "node", counter.Ident, "obj", extCount)
	counter.Merge(&extCount)
}

func (g *GCounter) IncVal(incr int) {
	g.Counter[g.Ident] += incr
}

func (g *GCounter) Count() (total int) {
	for _, v := range g.Counter {
		total += v
	}
	return
}

func (g *GCounter) Merge(c *GCounter) {
	for ident, val := range c.Counter {
		if v, ok := g.Counter[ident]; !ok || val > v {
			g.Counter[ident] = val
		}
	}
}

func main() {
	config := memberlist.DefaultLocalConfig()
	config.Delegate = &delegate{}

	counter.Ident = config.Name
	counter.Counter = make(map[string]int)

	http.HandleFunc("/inc", counter.incHandler)
	http.HandleFunc("/count", counter.countHandler)

	list, err := memberlist.Create(config)
	if err != nil {
		panic(err)
	}

	joinServer := os.Getenv("NODE")
	if joinServer != "localhost" {
		n, err := list.Join([]string{joinServer})
		if err != nil {
			panic("error joining node")
		}
		fmt.Printf("joined %v nodes\n", n)
	}

	go func() {
		for _, member := range list.Members() {
			fmt.Printf("Members: %s %s\n", member.Name, member.Addr)
		}
		time.Sleep(30 * time.Second)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	http.ListenAndServe(config.BindAddr+":3333", nil)

	<-stop
	if err := list.Leave(5 * time.Second); err != nil {
		panic(err)
	}

}

func (c *GCounter) incHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		c.IncVal(1)
		w.WriteHeader(204)
	}
}

func (c *GCounter) countHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		res := c.Count()
		json.NewEncoder(w).Encode(&res)
	}
}
