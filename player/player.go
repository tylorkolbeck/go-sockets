package player

import (
	"sync"

	"github.com/gorilla/websocket"
	math "github.com/tylorkolbeck/go-sockets/engine/Math"
)

type Player struct {
	ID   string
	Pos  math.Vec3
	Conn *websocket.Conn
	wMu  sync.Mutex
}

func NewPlayer(id string, conn *websocket.Conn, pos math.Vec3) *Player {
	return &Player{
		ID:   id,
		Pos:  pos,
		Conn: conn,
	}
}

func (p *Player) MoveUp() {
	p.wMu.Lock()
	defer p.wMu.Unlock()
	p.Pos.Y -= 1
}

func (p *Player) MoveDown() {
	p.wMu.Lock()
	defer p.wMu.Unlock()
	p.Pos.Y += 1
}

func (p *Player) MoveLeft() {
	p.wMu.Lock()
	defer p.wMu.Unlock()
	p.Pos.X -= 1
}

func (p *Player) MoveRight() {
	p.wMu.Lock()
	defer p.wMu.Unlock()
	p.Pos.X += 1
}
