package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	math "github.com/tylorkolbeck/go-sockets/engine/Math"
	"github.com/tylorkolbeck/go-sockets/player"
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

	p := &player.Player{
		ID:   id,
		Pos:  math.Vec3{X: 0, Y: 0, Z: 0},
		Conn: conn,
	}

	s.msgChannel <- JoinMsg{Type: "join", ID: id, Player: *p}

	go s.playerWsReadLoop(p)
}

func (s *Server) playerWsReadLoop(p *player.Player) {
	defer func() {
		s.msgChannel <- PlayerLeaveMsg{Type: "leave", ID: p.ID}
		if p.Conn != nil {
			_ = p.Conn.Close()
		}
	}()

	for {
		_, data, err := p.Conn.ReadMessage()
		if err != nil {
			return
		}

		var msg WsMsg
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		s.msgChannel <- WsMsg{
			ID:  msg.ID,
			Msg: msg.Msg,
		}
	}
}
