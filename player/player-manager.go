package player

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/tylorkolbeck/go-sockets/internal/logger"
	math "github.com/tylorkolbeck/go-sockets/lib/Math"
	"github.com/tylorkolbeck/go-sockets/models"
)

type PlayerManager struct {
	mu           sync.RWMutex
	players      map[string]*Player
	connToPlayer map[*websocket.Conn]*Player // Reverse lookup of player by conn
	logger       *logger.Logger
}

func NewPlayerManager() *PlayerManager {
	log := logger.NewLogger("Player Manager")
	log.Info("Initializing")
	return &PlayerManager{
		players:      make(map[string]*Player),
		connToPlayer: make(map[*websocket.Conn]*Player),
		logger:       log,
	}
}

// Player lifecycle
func (pm *PlayerManager) AddPlayer(id string, conn *websocket.Conn) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	p := NewPlayer(id, conn, math.Vec3{X: 0, Y: 0, Z: 0})
	pm.players[id] = p
	pm.connToPlayer[conn] = p
}

func (pm *PlayerManager) RemovePlayer(id string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	p, exists := pm.players[id]
	if exists {
		pm.logger.Info("Player left - ID: %s", p.ID)

		if p.Conn != nil {
			p.Conn.Close()
		}
		delete(pm.players, id)
		delete(pm.connToPlayer, p.Conn)
	}
}

// Player Data access (read-only)
func (pm *PlayerManager) GetPlayerPosition(id string) (*math.Vec3, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	p := pm.getPlayer(id)
	if p != nil {
		return &p.Pos, true
	} else {
		return nil, false
	}
}

func (pm *PlayerManager) GetPlayerConnection(id string) (*websocket.Conn, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	p := pm.getPlayer(id)
	if p != nil {
		return p.Conn, true
	} else {
		return nil, false
	}
}

func (pm *PlayerManager) GetAllPlayerSnapshots() map[string]models.PlayerSnapshot {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	playerPositions := make(map[string]models.PlayerSnapshot)

	for _, p := range pm.players {
		playerPositions[p.ID] = models.PlayerSnapshot{
			Pos: p.Pos,
		}

	}

	return playerPositions
}

// Player modifications
func (pm *PlayerManager) MovePlayer(id string, inputMapping InputMapping) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	p := pm.getPlayer(id)
	if p != nil {
		if inputMapping.Up {
			p.MoveUp()
		}
		if inputMapping.Down {
			p.MoveDown()
		}
		if inputMapping.Left {
			p.MoveLeft()

		}
		if inputMapping.Right {
			p.MoveRight()
		}
	}
}

// private
func (pm *PlayerManager) getPlayer(id string) *Player {
	return pm.players[id]
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

func (pm *PlayerManager) FindPlayerIDByConnection(conn *websocket.Conn) string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if player := pm.connToPlayer[conn]; player != nil {
		return player.ID
	}

	return ""
}
