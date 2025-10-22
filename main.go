package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tylorkolbeck/go-sockets/server"
)

var port = flag.String("port", "8000", "http service port")
var host = flag.String("host", "localhost", "http service host")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

func main() {
	fs := http.FileServer(http.Dir("frontend"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := server.NewServer()
	go server.Run(ctx)

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/connect", echo)
	http.Handle("/", fs)

	uri := fmt.Sprintf("%s:%s", *host, *port)
	log.Printf("Listening on: %s", uri)

	log.Fatal(http.ListenAndServe(uri, nil))

}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade: ", err)
		http.Error(w, err.Error(), 500)
		return
	}

	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read: ", err)
			break
		}

		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write: ", err)
			break
		}
	}
}
