package gameEngine

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tylorkolbeck/go-sockets/internal/config"
	"github.com/tylorkolbeck/go-sockets/internal/logger"
	"github.com/tylorkolbeck/go-sockets/models"
	"github.com/tylorkolbeck/go-sockets/player"
)

type GameEngine struct {
	msgChannel    chan any
	tick          uint64
	worldW        float64
	worldH        float64
	worldbg       models.Color
	tickRate      int
	playerManager *player.PlayerManager
	maxPlayers    int
	logger        *logger.Logger
}

func NewGameEngine(gameConfig config.GameConfig) *GameEngine {
	log := logger.NewLogger("Game Engine")
	log.Info("Initializing")

	msgChannel := make(chan any, 1024)
	return &GameEngine{
		msgChannel: msgChannel,
		worldW:     gameConfig.WorldWidth,
		worldH:     gameConfig.WorldHeight,
		worldbg: models.Color{
			R: 255,
			G: 230,
			B: 0,
		},
		tick:          0,
		tickRate:      gameConfig.TickRate,
		maxPlayers:    gameConfig.MaxPlayers,
		playerManager: player.NewPlayerManager(),
		logger:        log,
	}
}

func (ge *GameEngine) GameLoop(ctx context.Context) {
	if ge.tick == 0 {
		ge.logger.Info("Game Loop Started")
	}
	ticker := time.NewTicker(time.Duration(ge.tickRate) * time.Millisecond) // 20hz
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-ge.msgChannel:
			switch e := ev.(type) {
			case player.JoinMsg:
				ge.playerManager.AddPlayer(e.ID, e.Conn)
				ge.broadcastWorldSettings(e.ID)
				// Need to tell everyone a player joined
				ge.broadcastPlayerList()
				ge.broadcastPlayerJoined(e.ID)
				ge.logger.Info("Player joined - ID: %s", e.ID)
			case player.PlayerLeaveMsg:
				ge.playerManager.RemovePlayer(e.ID)
				ge.broadcastPlayerLeft(e.ID)
			case player.RawWsMsg:
				ge.playerManager.MovePlayer(e.ID, e.Data)
			}
		case <-ticker.C:
			ge.tick++
			ge.broadcastPlayerSnapshots()
		}
	}
}

func (ge *GameEngine) OnMessageHandler(conn *websocket.Conn, msgType int, data []byte) error {
	var msg player.PlayerWsMsg
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	ge.msgChannel <- player.RawWsMsg{
		ID:   msg.ID,
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

	ge.msgChannel <- player.JoinMsg{Type: "join", ID: id, Conn: conn}

	// ge.playerManager.AddPlayer(id, conn)
	return nil
}

func (ge *GameEngine) OnClientDisconnectHandler(conn *websocket.Conn) {
	playerId := ge.playerManager.FindPlayerIDByConnection(conn)
	if playerId != "" {
		ge.msgChannel <- player.PlayerLeaveMsg{Type: "leave", ID: playerId}
	}
}
