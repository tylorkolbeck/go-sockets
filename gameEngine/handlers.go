package gameEngine

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/tylorkolbeck/go-sockets/player"
)

func (gm *GameEngine) broadcastWorldSettings(player player.Player) {
	worldSettingsMsg := WorldSettingsMsg{
		Type:    "worldsettings",
		WorldW:  gm.worldW,
		WorldH:  gm.worldH,
		WorldBg: gm.worldbg,
	}

	data, _ := json.Marshal(worldSettingsMsg)
	if player.Conn != nil {
		player.Conn.WriteMessage(websocket.TextMessage, data)
	}
}

func (gm *GameEngine) broadcastPlayerLeft(id string) {
	leftPlayerMsg := player.PlayerLeaveMsg{
		Type: "playerleft",
		ID:   id,
	}
	data, _ := json.Marshal(leftPlayerMsg)
	for _, p := range gm.playerManager.GetAllPlayers() {
		if p.Conn != nil {
			p.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}

func (gm *GameEngine) broadcastPlayerJoined(id string) {
	joinedPlayerMsg := PlayerJoinMsg{
		Type: "playerjoined",
		ID:   id,
	}
	data, _ := json.Marshal(joinedPlayerMsg)
	for _, p := range gm.playerManager.GetAllPlayers() {
		if p.ID != id && p.Conn != nil {
			p.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}

func (gm *GameEngine) broadcastPlayerList() {
	var playerIds []string
	players := gm.playerManager.GetAllPlayers()

	for _, p := range players {
		playerIds = append(playerIds, p.ID)
	}

	for _, p := range players {
		if p.Conn != nil {
			updatedPlayerList := PlayerListMsg{
				Type:      "updatedplayerlist",
				PlayerIds: playerIds,
			}
			data, _ := json.Marshal(updatedPlayerList)
			_ = p.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}

func (gm *GameEngine) broadcast() {
	snapshot := SnapshotMsg{
		Type:    "snapshot",
		Tick:    gm.tick,
		Players: map[string]PlayerSnapshot{},
	}

	for id, p := range gm.playerManager.GetAllPlayers() {
		snapshot.Players[id] = PlayerSnapshot{
			Pos: p.Pos,
		}
	}

	data, _ := json.Marshal(snapshot)
	for _, p := range gm.playerManager.GetAllPlayers() {
		if p.Conn != nil {
			_ = p.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}
