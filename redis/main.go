package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client
var sub *redis.PubSub
var Users = map[string]*websocket.Conn{}

const channelName = "chat"

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("ping failed. could not connect", err)
	}

	go func() {
		sub = client.Subscribe(context.Background(), channelName)
		ch := sub.Channel()

		for message := range ch {
			s := strings.Split(message.Payload, ":")
			for user, peer := range Users {
				if user != s[0] {
					msg := "[" + s[0] + " says]: " + s[1]
					peer.WriteMessage(1, []byte(msg))
				}
			}
			log.Println(message.Payload)
		}
	}()

	http.HandleFunc("/chat/", wsConnect)
	server := http.Server{Addr: ":8080", Handler: nil}

	go func() {
		log.Println("chat server started")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, conn := range Users {
		conn.Close()
	}

	sub.Unsubscribe(context.Background(), channelName)
	sub.Close()

	server.Shutdown(ctx)

	fmt.Println("chat application closed")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsConnect(w http.ResponseWriter, req *http.Request) {
	name := strings.TrimPrefix(req.URL.Path, "/chat/")

	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("protocol upgrade error", err)
		return

	}
	Users[name] = c
	log.Println(name, "joined the chat")

	for {
		_, p, err := c.ReadMessage()
		if err != nil {
			_, ok := err.(*websocket.CloseError)
			if ok {
				fmt.Println("connection closed by:", name)
				err := c.Close()
				if err != nil {
					fmt.Println("close connection error", err)
				}
				delete(Users, name)
				fmt.Println("connection and user session closed")
			}
			break
		}
		err = client.Publish(context.Background(), "chat", name+":"+string(p)).Err()
		if err != nil {
			log.Println("publish error", err)
		}
	}
}
