package server

import (
	math "github.com/tylorkolbeck/go-sockets/engine/Math"
	"github.com/tylorkolbeck/go-sockets/player"
)

type JoinMsg struct {
	Type   string        `json:"type"`
	ID     string        `json:"id"`
	Player player.Player `json:"player"`
}

type PlayerJoinMsg struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type PlayerListMsg struct {
	Type      string   `json:"type"`
	PlayerIds []string `json:"playerIds"`
}

type PlayerLeaveMsg struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type InputMsg struct {
	Type  string `json:"type"`
	Up    bool   `json:"up"`
	Down  bool   `json:"down"`
	Left  bool   `json:"left"`
	Right bool   `json:"right"`
}

type WsMsg struct {
	ID  string   `json:"id"`
	Msg InputMsg `json:"msg"`
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
	Type    string  `json:"type"`
	WorldW  float64 `json:"worldwidth"`
	WorldH  float64 `json:"worldheight"`
	WorldBg Color   `json:"worldbg"`
}
