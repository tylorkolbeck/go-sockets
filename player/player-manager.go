package player

import (
	"log"

	"github.com/gorilla/websocket"
	math "github.com/tylorkolbeck/go-sockets/engine/Math"
)

type PlayerManager struct {
	players    map[string]*Player
	msgChannel chan any
}

func NewPlayerManager(msgChannel chan any) *PlayerManager {
	return &PlayerManager{
		players:    make(map[string]*Player),
		msgChannel: msgChannel,
	}
}

func (pm *PlayerManager) GetPlayer(id string) *Player {
	return pm.players[id]
}

func (pm *PlayerManager) AddPlayer(id string, conn *websocket.Conn) {
	p := NewPlayer(id, math.Vec3{X: 0, Y: 0, Z: 0}, conn)
	pm.players[id] = p

	pm.msgChannel <- JoinMsg{Type: "join", ID: id, Player: *p}
	go p.StartPlayerWsReadLoop(pm.msgChannel)
}

func (pm *PlayerManager) GetAllPlayers() map[string]*Player {
	return pm.players
}

func (pm *PlayerManager) RemovePlayer(id string) {
	if p, ok := pm.players[id]; ok {
		log.Printf("Player left - ID: %s", p.ID)
		if p.Conn != nil {
			_ = p.Conn.Close()
		}

		delete(pm.players, id)
	}
}
