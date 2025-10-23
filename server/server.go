package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	math "github.com/tylorkolbeck/go-sockets/engine/Math"
	"github.com/tylorkolbeck/go-sockets/player"
)

type JoinMsg struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type LeaveMsg struct {
	Type string `json:"leave"`
	ID   string `json:"id"`
}

type InputMsg struct {
	Type  string `json:"type"`
	Up    bool   `json:"up"`
	Down  bool   `json:"down"`
	Left  bool   `json:"left"`
	Right bool   `json:"right"`
}

type WsMsg struct {
	ID  string   `json:"id"`
	Msg InputMsg `json:"msg"`
}

type Server struct {
	mu         sync.Mutex
	players    map[string]*player.Player
	msgChannel chan any
	tick       uint64
	worldW     float64
	worldH     float64
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

func NewServer() *Server {
	return &Server{
		players:    map[string]*player.Player{},
		msgChannel: make(chan any, 1024),
		worldW:     800,
		worldH:     800,
		tick:       0,
	}
}

func (s *Server) Run(ctx context.Context) {
	ticker := time.NewTicker(50 * time.Millisecond) // 20hz
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.tick++
		case ev := <-s.msgChannel:
			switch e := ev.(type) {
			case JoinMsg:
				s.mu.Lock()
				s.players[e.ID] = &player.Player{
					ID:  e.ID,
					Pos: math.Vec3{X: 100, Y: 100, Z: 0},
				}
				log.Printf("Player joined - ID: %s", e.ID)
				s.mu.Unlock()
			case LeaveMsg:
				s.mu.Lock()
				// Check if player exists
				if p, ok := s.players[e.ID]; ok {
					log.Printf("Player left - ID: %s", p.ID)
					if p.Conn != nil {
						_ = p.Conn.Close()
					}
					delete(s.players, e.ID)
				}
				s.mu.Unlock()
			case WsMsg:
				s.mu.Lock()
				p, _ := s.GetPlayer(e.ID)
				if p != nil {
					log.Printf("Input Message - ID: %s", p.ID)
				}

				s.mu.Unlock()
			}
		}
	}
}

func (s *Server) HandleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), 600)
		return
	}

	id := r.URL.Query().Get("id")
	s.msgChannel <- JoinMsg{Type: "join", ID: id}

	p := &player.Player{
		ID:   id,
		Pos:  math.Vec3{X: 0, Y: 0, Z: 0},
		Conn: conn,
	}

	go s.playerReadLoop(p)
}

func (s *Server) playerReadLoop(p *player.Player) {
	defer func() {
		s.msgChannel <- LeaveMsg{Type: "leave", ID: p.ID}
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

		log.Printf("%+v", msg)

	}
}

func (s *Server) GetPlayer(id string) (*player.Player, bool) {
	player, exists := s.players[id]

	return player, exists
}
