package player

import (
	"github.com/gorilla/websocket"
	math "github.com/tylorkolbeck/go-sockets/engine/Math"
)

type Player struct {
	ID   string
	Pos  math.Vec3
	Conn *websocket.Conn
}
