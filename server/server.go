package server

import (
	"context"
	"time"

	"github.com/tylorkolbeck/go-sockets/player"
)

type Server struct {
	players map[string]*player.Player
	tick    uint64
	worldW  float64
	worldH  float64
}

func NewServer() *Server {
	return &Server{
		players: map[string]*player.Player{},
		worldW:  800,
		worldH:  800,
		tick:    0,
	}
}

func (s *Server) Run(ctx context.Context) {
	ticker := time.NewTicker(50 * time.Millisecond) // 20hz
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.tick++
		}
	}

}
