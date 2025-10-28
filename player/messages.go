package player

import "github.com/gorilla/websocket"

type JoinMsg struct {
	Type string          `json:"type"`
	ID   string          `json:"id"`
	Conn *websocket.Conn `json:"conn"`
}

type PlayerLeaveMsg struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type PlayerWsMsg struct {
	ID  string       `json:"id"`
	Msg InputMapping `json:"msg"`
}

type RawWsMsg struct {
	ID   string
	Conn *websocket.Conn
	Data InputMapping
}
