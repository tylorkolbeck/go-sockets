package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

func (s *Server) HandleWs(w http.ResponseWriter, r *http.Request) {
	// Check for required parameters BEFORE upgrading connection
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("Connection attempted without a user id. Not upgrading connection")
		http.Error(w, "Missing id parameter", 400)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	s.playerManager.AddPlayer(id, conn)
}
