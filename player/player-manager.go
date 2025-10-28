package player

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
	math "github.com/tylorkolbeck/go-sockets/engine/Math"
)

type PlayerManager struct {
	mu           sync.RWMutex
	players      map[string]*Player
	connToPlayer map[*websocket.Conn]*Player // Reverse lookup of player by conn
}

func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		players:      make(map[string]*Player),
		connToPlayer: make(map[*websocket.Conn]*Player),
	}
}

func (pm *PlayerManager) FindPlayerByConnection(conn *websocket.Conn) *Player {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.connToPlayer[conn]
}

func (pm *PlayerManager) GetPlayer(id string) *Player {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.players[id]
}

func (pm *PlayerManager) AddPlayer(id string, conn *websocket.Conn) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	p := NewPlayer(id, conn, math.Vec3{X: 0, Y: 0, Z: 0})
	pm.players[id] = p
	pm.connToPlayer[conn] = p
}

func (pm *PlayerManager) GetAllPlayers() map[string]*Player {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	playersCopy := make(map[string]*Player, len(pm.players))
	for k, v := range pm.players {
		playersCopy[k] = v
	}

	return playersCopy
}

func (pm *PlayerManager) RemovePlayer(id string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	p, exists := pm.players[id]
	if exists {
		log.Printf("Player left - ID: %s", p.ID)

		if p.Conn != nil {
			p.Conn.Close()
		}
		delete(pm.players, id)
		delete(pm.connToPlayer, p.Conn)
	}
}

// Utilities
func (pm *PlayerManager) GetPlayerIDs() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	ids := make([]string, 0, len(pm.players))
	for id := range pm.players {
		ids = append(ids, id)
	}

	return ids
}

func (pm *PlayerManager) GetPlayerCount() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return len(pm.players)
}

func (pm *PlayerManager) PlayerExists(id string) bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	_, exists := pm.players[id]
	return exists
}
