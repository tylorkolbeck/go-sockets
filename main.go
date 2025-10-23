package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tylorkolbeck/go-sockets/server"
)

var port = flag.String("port", "8000", "http service port")
var host = flag.String("host", "localhost", "http service host")

func main() {
	fs := http.FileServer(http.Dir("frontend"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, stopping server")
		cancel()

		os.Exit(0)
	}()

	server := server.NewServer()
	go server.Run(ctx)

	flag.Parse()
	http.HandleFunc("/connect", server.HandleWs)
	http.Handle("/", fs)

	uri := fmt.Sprintf("%s:%s", *host, *port)
	log.Printf("Listening on: %s", uri)

	log.Fatal(http.ListenAndServe(uri, nil))

}
