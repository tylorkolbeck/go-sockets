package server

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/tylorkolbeck/go-sockets/models"
	"github.com/tylorkolbeck/go-sockets/player"
)

type Server struct {
	mu            sync.Mutex
	msgChannel    chan any
	tick          uint64
	worldW        float64
	worldH        float64
	worldbg       models.Color
	playerManager *player.PlayerManager
}

func NewServer() *Server {
	msgChannel := make(chan any, 1024)
	return &Server{
		msgChannel: msgChannel,
		worldW:     800,
		worldH:     800,
		worldbg: models.Color{
			R: 255,
			G: 230,
			B: 0,
		},
		tick:          0,
		playerManager: player.NewPlayerManager(msgChannel),
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
			case player.JoinMsg:
				s.mu.Lock()
				s.broadcastWorldSettings(e.Player)

				// Need to tell everyone a player joined
				s.broadcastPlayerList()
				s.broadcastPlayerJoined(e.ID)
				log.Printf("Player joined - ID: %s", e.ID)

				s.mu.Unlock()
			case models.PlayerLeaveMsg:
				s.mu.Lock()
				s.playerManager.RemovePlayer(e.ID)
				s.broadcastPlayerLeft(e.ID)
				s.mu.Unlock()
			case player.PlayerWsMsg:
				s.mu.Lock()
				p := s.playerManager.GetPlayer(e.ID)
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
