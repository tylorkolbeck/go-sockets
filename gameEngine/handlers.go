package gameEngine

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/tylorkolbeck/go-sockets/models"
	"github.com/tylorkolbeck/go-sockets/player"
)

// Single client broadcasts
func (gm *GameEngine) broadcastWorldSettings(id string) {
	worldSettingsMsg := models.WorldSettingsMsg{
		Type:    models.MsgTypeWorldSettings,
		WorldW:  gm.worldW,
		WorldH:  gm.worldH,
		WorldBg: gm.worldbg,
	}

	data, err := json.Marshal(worldSettingsMsg)
	if err != nil {
		log.Printf("Failed to marshal snapshot: %v", err)
	}

	conn, ok := gm.playerManager.GetPlayerConnection(id)
	if ok {
		broadcastTextMsg(conn, data)
	}
}

// All client broadcasts
func (gm *GameEngine) broadcastPlayerLeft(id string) {
	leftPlayerMsg := player.PlayerLeaveMsg{
		Type: models.MsgTypePlayerLeft,
		ID:   id,
	}
	data, err := json.Marshal(leftPlayerMsg)
	if err != nil {
		log.Printf("Failed to marshal snapshot: %v", err)
	}

	gm.broadCastToAllPlayers(data)
}

func (gm *GameEngine) broadcastPlayerJoined(id string) {
	joinedPlayerMsg := models.PlayerJoinMsg{
		Type: models.MsgTypePlayerJoined,
		ID:   id,
	}
	data, err := json.Marshal(joinedPlayerMsg)
	if err != nil {
		log.Printf("Failed to marshal snapshot: %v", err)
	}

	gm.broadCastToAllPlayers(data)
}

func (gm *GameEngine) broadcastPlayerList() {
	var playerIds = gm.playerManager.GetPlayerIDs()

	playerListMsg := models.PlayerListMsg{
		Type:      models.MsgTypePlayerList,
		PlayerIds: playerIds,
	}

	data, err := json.Marshal(playerListMsg)
	if err != nil {
		log.Printf("Failed to marshal snapshot: %v", err)
	}

	gm.broadCastToAllPlayers(data)
}

func (gm *GameEngine) broadcastPlayerSnapshots() {
	snapshot := models.SnapshotMsg{
		Type:    models.MsgTypeSnapshot,
		Tick:    gm.tick,
		Players: gm.playerManager.GetAllPlayerSnapshots(),
	}

	data, err := json.Marshal(snapshot)
	if err != nil {
		log.Printf("Failed to marshal snapshot: %v", err)
	}

	gm.broadCastToAllPlayers(data)
}

// Utilities
func (gm *GameEngine) broadCastToAllPlayers(data []byte) {
	playerIds := gm.playerManager.GetPlayerIDs()
	for _, id := range playerIds {
		conn, ok := gm.playerManager.GetPlayerConnection(id)
		if ok {
			broadcastTextMsg(conn, data)
		}
	}
}

func broadcastTextMsg(conn *websocket.Conn, data []byte) {
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
