package server

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/tylorkolbeck/go-sockets/player"
)

type Color struct {
	R int64 `json:"r"`
	G int64 `json:"g"`
	B int64 `json:"b"`
}

type Server struct {
	mu         sync.Mutex
	players    map[string]*player.Player
	msgChannel chan any
	tick       uint64
	worldW     float64
	worldH     float64
	worldbg    Color
}

func NewServer() *Server {
	return &Server{
		players:    map[string]*player.Player{},
		msgChannel: make(chan any, 1024),
		worldW:     800,
		worldH:     800,
		worldbg: Color{
			R: 255,
			G: 230,
			B: 0,
		},
		tick: 0,
	}
}

func (s *Server) update() {
	s.mu.Lock()
	defer s.mu.Unlock()
}

func (s *Server) Run(ctx context.Context) {
	ticker := time.NewTicker(50 * time.Millisecond) // 20hz
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-s.msgChannel:
			switch e := ev.(type) {
			case JoinMsg:
				s.mu.Lock()
				s.players[e.ID] = &e.Player

				s.broadcastWorldSettings(e.Player)

				// Need to tell everyone a player joined
				s.broadcastPlayerList()
				s.broadcastPlayerJoined(e.ID)
				log.Printf("Player joined - ID: %s", e.ID)

				s.mu.Unlock()
			case PlayerLeaveMsg:
				s.mu.Lock()
				if p, ok := s.players[e.ID]; ok {
					log.Printf("Player left - ID: %s", p.ID)
					if p.Conn != nil {
						_ = p.Conn.Close()
					}
					s.broadcastPlayerLeft(e.ID)
					delete(s.players, e.ID)
				}
				s.mu.Unlock()
			case WsMsg:
				s.mu.Lock()
				p, _ := s.GetPlayer(e.ID)
				if p != nil {
					if e.Msg.Up {
						p.MoveUp()
					}
					if e.Msg.Down {
						p.MoveDown()
					}
					if e.Msg.Left {
						p.MoveLeft()

					}
					if e.Msg.Right {
						p.MoveRight()
					}
				}

				s.mu.Unlock()
			}
		case <-ticker.C:
			s.update()
			s.tick++
			s.broadcast()
		}
	}
}

func (s *Server) GetPlayer(id string) (*player.Player, bool) {
	player, exists := s.players[id]

	return player, exists
}
