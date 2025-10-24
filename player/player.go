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

func (p *Player) MoveUp() {
	p.Pos.Y -= 1
}

func (p *Player) MoveDown() {
	p.Pos.Y += 1
}

func (p *Player) MoveLeft() {
	p.Pos.X -= 1
}

func (p *Player) MoveRight() {
	p.Pos.X += 1
}
