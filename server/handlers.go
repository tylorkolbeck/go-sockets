package server

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/tylorkolbeck/go-sockets/player"
)

func (s *Server) broadcastWorldSettings(player player.Player) {
	worldSettingsMsg := WorldSettingsMsg{
		Type:    "worldsettings",
		WorldW:  s.worldW,
		WorldH:  s.worldH,
		WorldBg: s.worldbg,
	}

	data, _ := json.Marshal(worldSettingsMsg)
	if player.Conn != nil {
		player.Conn.WriteMessage(websocket.TextMessage, data)
	}
}

func (s *Server) broadcastPlayerLeft(id string) {
	leftPlayerMsg := player.PlayerLeaveMsg{
		Type: "playerleft",
		ID:   id,
	}
	data, _ := json.Marshal(leftPlayerMsg)
	for _, p := range s.playerManager.GetAllPlayers() {
		if p.Conn != nil {
			p.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}

func (s *Server) broadcastPlayerJoined(id string) {
	joinedPlayerMsg := PlayerJoinMsg{
		Type: "playerjoined",
		ID:   id,
	}
	data, _ := json.Marshal(joinedPlayerMsg)
	for _, p := range s.playerManager.GetAllPlayers() {
		if p.ID != id && p.Conn != nil {
			p.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}

func (s *Server) broadcastPlayerList() {
	var playerIds []string
	players := s.playerManager.GetAllPlayers()

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

func (s *Server) broadcast() {
	s.mu.Lock()
	defer s.mu.Unlock()

	snapshot := SnapshotMsg{
		Type:    "snapshot",
		Tick:    s.tick,
		Players: map[string]PlayerSnapshot{},
	}

	for id, p := range s.playerManager.GetAllPlayers() {
		snapshot.Players[id] = PlayerSnapshot{
			Pos: p.Pos,
		}
	}

	data, _ := json.Marshal(snapshot)
	for _, p := range s.playerManager.GetAllPlayers() {
		if p.Conn != nil {
			_ = p.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}
