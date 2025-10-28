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
	"time"

	"github.com/tylorkolbeck/go-sockets/gameEngine"
	"github.com/tylorkolbeck/go-sockets/internal/config"
	"github.com/tylorkolbeck/go-sockets/internal/logger"
	wsm "github.com/tylorkolbeck/go-sockets/websocket"
)

var port = flag.String("port", "8000", "http service port")
var host = flag.String("host", "localhost", "http service host")

func main() {
	logger := logger.NewLogger("main")
	fs := http.FileServer(http.Dir("frontend"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Received shutdown signal, stopping server")
		cancel()

		time.Sleep(100 * time.Millisecond)
		os.Exit(0)
	}()

	config := config.Config{
		Game: config.GameConfig{
			TickRate:    20,
			WorldWidth:  800,
			WorldHeight: 800,
			MaxPlayers:  5,
		},
		Server: config.ServerConfig{
			Host: *host,
			Port: *port,
		},
	}

	gameEngine := gameEngine.NewGameEngine(config.Game)
	wsManager := wsm.NewWebsocketManager(gameEngine.OnMessageHandler, gameEngine.OnClientConnectHandler, gameEngine.OnClientDisconnectHandler)

	go gameEngine.GameLoop(ctx)

	// Websocket route
	http.HandleFunc("/connect", wsManager.HandleConnection)

	// File Server route
	http.Handle("/", fs)

	uri := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	logger.Info("Listening on: %s", uri)

	log.Fatal(http.ListenAndServe(uri, nil))
}
