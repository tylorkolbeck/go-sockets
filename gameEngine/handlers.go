package gameEngine

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/tylorkolbeck/go-sockets/models"
	"github.com/tylorkolbeck/go-sockets/player"
)

// Single client broadcasts
func (ge *GameEngine) broadcastWorldSettings(id string) {
	worldSettingsMsg := models.WorldSettingsMsg{
		Type:    models.MsgTypeWorldSettings,
		WorldW:  ge.worldW,
		WorldH:  ge.worldH,
		WorldBg: ge.worldbg,
	}

	data, err := json.Marshal(worldSettingsMsg)
	if err != nil {
		ge.logger.Error("Failed to marshal snapshot: %v", err)
	}

	conn, ok := ge.playerManager.GetPlayerConnection(id)
	if ok {
		ge.broadcastTextMsg(conn, data)
	}
}

// All client broadcasts
func (ge *GameEngine) broadcastPlayerLeft(id string) {
	leftPlayerMsg := player.PlayerLeaveMsg{
		Type: models.MsgTypePlayerLeft,
		ID:   id,
	}
	data, err := json.Marshal(leftPlayerMsg)
	if err != nil {
		ge.logger.Error("Failed to marshal snapshot: %v", err)
	}

	ge.broadCastToAllPlayers(data)
}

func (ge *GameEngine) broadcastPlayerJoined(id string) {
	joinedPlayerMsg := models.PlayerJoinMsg{
		Type: models.MsgTypePlayerJoined,
		ID:   id,
	}
	data, err := json.Marshal(joinedPlayerMsg)
	if err != nil {
		ge.logger.Error("Failed to marshal snapshot: %v", err)
	}

	ge.broadCastToAllPlayers(data)
}

func (ge *GameEngine) broadcastPlayerList() {
	var playerIds = ge.playerManager.GetPlayerIDs()

	playerListMsg := models.PlayerListMsg{
		Type:      models.MsgTypePlayerList,
		PlayerIds: playerIds,
	}

	data, err := json.Marshal(playerListMsg)
	if err != nil {
		ge.logger.Error("Failed to marshal snapshot: %v", err)
	}

	ge.broadCastToAllPlayers(data)
}

func (ge *GameEngine) broadcastPlayerSnapshots() {
	snapshot := models.SnapshotMsg{
		Type:    models.MsgTypeSnapshot,
		Tick:    ge.tick,
		Players: ge.playerManager.GetAllPlayerSnapshots(),
	}

	data, err := json.Marshal(snapshot)
	if err != nil {
		ge.logger.Error("Failed to marshal snapshot: %v", err)

	}

	ge.broadCastToAllPlayers(data)
}

// Utilities
func (ge *GameEngine) broadCastToAllPlayers(data []byte) {
	playerIds := ge.playerManager.GetPlayerIDs()
	for _, id := range playerIds {
		conn, ok := ge.playerManager.GetPlayerConnection(id)
		if ok {
			ge.broadcastTextMsg(conn, data)
		}
	}
}

func (ge *GameEngine) broadcastTextMsg(conn *websocket.Conn, data []byte) {
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		ge.logger.Error("Failed to send message: %v", err)
	}
}
