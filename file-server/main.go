package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/pkg/browser"
)

var (
	path     = flag.String("d", ".", "path to the folder to serve. Default current folder")
	port     = flag.String("p", "8080", "port to serve on. Default 8080")
	autoOpen = flag.Bool("auto", false, "auto open browser. Default disabled")
)

func main() {
	flag.Parse()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	dirname, err := filepath.Abs(*path)
	if err != nil {
		log.Fatalf("Could not get absolute path to directory: %s: %s", dirname, err.Error())
	}

	httpServer := http.Server{
		Addr: ":" + *port,
	}

	go func() {
		log.Printf("Serving %s on port %s", dirname, *port)

		fs := http.FileServer(http.Dir(dirname))
		mux := http.NewServeMux()
		mux.Handle("/", fs)
		httpServer.Handler = mux

		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe Error: %v", err)
		}
	}()

	if *autoOpen {
		time.Sleep(time.Second * 2)
		browser.OpenURL("http://localhost:" + *port)
	}

	<-sigint
	if err := httpServer.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP Server Shutdown Error: %v", err)
	}

	log.Println("Server stop")

}
