package gameEngine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tylorkolbeck/go-sockets/models"
	"github.com/tylorkolbeck/go-sockets/player"
)

type GameEngine struct {
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

func (ge *GameEngine) StartGameLoop(ctx context.Context) {
	ticker := time.NewTicker(50 * time.Millisecond) // 20hz
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-ge.msgChannel:
			switch e := ev.(type) {
			case player.JoinMsg:
				ge.broadcastWorldSettings(ge.playerManager.GetPlayer(e.ID))

				// Need to tell everyone a player joined
				ge.broadcastPlayerList()
				ge.broadcastPlayerJoined(e.ID)
				log.Printf("Player joined - ID: %s", e.ID)
			case player.PlayerLeaveMsg:
				ge.playerManager.RemovePlayer(e.ID)
				ge.broadcastPlayerLeft(e.ID)
			case player.RawWsMsg:
				p := ge.playerManager.FindPlayerByConnection(e.Conn)
				if p != nil {
					if e.Data.Up {
						p.MoveUp()
					}
					if e.Data.Down {
						p.MoveDown()
					}
					if e.Data.Left {
						p.MoveLeft()

					}
					if e.Data.Right {
						p.MoveRight()
					}
					// ge.msgChannel <- player.PlayerWsMsg{
					// 	ID:  p.ID,
					// 	Msg: e.Data,
					// }
				}

			}
		case <-ticker.C:
			ge.tick++
			ge.broadcast()
		}
	}
}

func (ge *GameEngine) OnMessageHandler(conn *websocket.Conn, msgType int, data []byte) error {
	var msg player.PlayerWsMsg
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	ge.msgChannel <- player.RawWsMsg{
		Conn: conn,
		Data: msg.Msg,
	}

	return nil
}

func (ge *GameEngine) OnClientConnectHandler(conn *websocket.Conn, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		return fmt.Errorf("missing id parameter")
	}

	ge.playerManager.AddPlayer(id, conn)
	return nil
}

func (ge *GameEngine) OnClientDisconnectHandler(conn *websocket.Conn) {
	p := ge.playerManager.FindPlayerByConnection(conn)
	if p != nil {
		ge.msgChannel <- player.PlayerLeaveMsg{Type: "leave", ID: p.ID}
	}

}
