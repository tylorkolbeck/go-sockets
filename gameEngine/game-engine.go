package gameEngine

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/tylorkolbeck/go-sockets/models"
	"github.com/tylorkolbeck/go-sockets/player"
)

type GameEngine struct {
	mu            sync.Mutex
	msgChannel    chan any
	tick          uint64
	worldW        float64
	worldH        float64
	worldbg       models.Color
	playerManager *player.PlayerManager
}

func NewGameEngine() *GameEngine {
	msgChannel := make(chan any, 1024)
	return &GameEngine{
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

// func (s *Server) update() {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// }

func (ge *GameEngine) Run(ctx context.Context) {
	ticker := time.NewTicker(50 * time.Millisecond) // 20hz
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-ge.msgChannel:
			switch e := ev.(type) {
			case player.JoinMsg:
				ge.mu.Lock()
				ge.broadcastWorldSettings(e.Player)

				// Need to tell everyone a player joined
				ge.broadcastPlayerList()
				ge.broadcastPlayerJoined(e.ID)
				log.Printf("Player joined - ID: %s", e.ID)

				ge.mu.Unlock()
			case models.PlayerLeaveMsg:
				ge.mu.Lock()
				ge.playerManager.RemovePlayer(e.ID)
				ge.broadcastPlayerLeft(e.ID)
				ge.mu.Unlock()
			case player.PlayerWsMsg:
				ge.mu.Lock()
				p := ge.playerManager.GetPlayer(e.ID)
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

				ge.mu.Unlock()
			}
		case <-ticker.C:
			// s.update()
			ge.tick++
			ge.broadcast()
		}
	}
}
