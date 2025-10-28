package player

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	math "github.com/tylorkolbeck/go-sockets/engine/Math"
	"github.com/tylorkolbeck/go-sockets/models"
)

type Player struct {
	ID   string
	Pos  math.Vec3
	Conn *websocket.Conn
}

func NewPlayer(id string, pos math.Vec3, conn *websocket.Conn) *Player {
	return &Player{
		ID:   id,
		Pos:  pos,
		Conn: conn,
	}
}

func (p *Player) StartPlayerWsReadLoop(msgChannel chan any) {
	defer func() {
		msgChannel <- models.PlayerLeaveMsg{Type: "leave", ID: p.ID}
		if p.Conn != nil {
			_ = p.Conn.Close()
		}
	}()

	for {
		_, data, err := p.Conn.ReadMessage()
		if err != nil {
			return
		}

		var msg PlayerWsMsg
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		msgChannel <- PlayerWsMsg{
			ID:  msg.ID,
			Msg: msg.Msg,
		}
	}
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
