package server

import (
	math "github.com/tylorkolbeck/go-sockets/engine/Math"
	"github.com/tylorkolbeck/go-sockets/models"
)

type PlayerJoinMsg struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type PlayerListMsg struct {
	Type      string   `json:"type"`
	PlayerIds []string `json:"playerIds"`
}

type PlayerSnapshot struct {
	Pos math.Vec3 `json:"pos"`
}

type SnapshotMsg struct {
	Type    string                    `json:"type"`
	Tick    uint64                    `json:"tick"`
	Players map[string]PlayerSnapshot `json:"players"`
}

type WorldSettingsMsg struct {
	Type    string       `json:"type"`
	WorldW  float64      `json:"worldwidth"`
	WorldH  float64      `json:"worldheight"`
	WorldBg models.Color `json:"worldbg"`
}
