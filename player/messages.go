package player

import "github.com/gorilla/websocket"

type JoinMsg struct {
	Type string `json:"type"`
	ID   string `json:"id"`
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
	Conn *websocket.Conn
	Data InputMapping
}
