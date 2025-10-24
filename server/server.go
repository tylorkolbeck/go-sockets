package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	math "github.com/tylorkolbeck/go-sockets/engine/Math"
	"github.com/tylorkolbeck/go-sockets/player"
)

type Color struct {
	R int64 `json:"r"`
	G int64 `json:"g"`
	B int64 `json:"b"`
}

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

type Server struct {
	mu         sync.Mutex
	players    map[string]*player.Player
	msgChannel chan any
	tick       uint64
	worldW     float64
	worldH     float64
	worldbg    Color
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

func NewServer() *Server {
	return &Server{
		players:    map[string]*player.Player{},
		msgChannel: make(chan any, 1024),
		worldW:     800,
		worldH:     800,
		worldbg: Color{
			R: 255,
			G: 230,
			B: 0,
		},
		tick: 0,
	}
}

func (s *Server) update() {
	s.mu.Lock()
	defer s.mu.Unlock()
}

func (s *Server) Run(ctx context.Context) {
	ticker := time.NewTicker(50 * time.Millisecond) // 20hz
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-s.msgChannel:
			switch e := ev.(type) {
			case JoinMsg:
				s.mu.Lock()
				s.players[e.ID] = &e.Player

				s.broadcastWorldSettings(e.Player)

				// Need to tell everyone a player joined
				s.broadcastPlayerList()
				s.broadcastPlayerJoined(e.ID)
				log.Printf("Player joined - ID: %s", e.ID)

				s.mu.Unlock()
			case PlayerLeaveMsg:
				s.mu.Lock()
				if p, ok := s.players[e.ID]; ok {
					log.Printf("Player left - ID: %s", p.ID)
					if p.Conn != nil {
						_ = p.Conn.Close()
					}
					s.broadcastPlayerLeft(e.ID)
					delete(s.players, e.ID)
				}
				s.mu.Unlock()
			case WsMsg:
				s.mu.Lock()
				p, _ := s.GetPlayer(e.ID)
				if p != nil {
					if e.Msg.Up {
						p.MoveUp()
					}
					if e.Msg.Down {
						p.MoveDown()
					}
					if e.Msg.Left {
						p.MoveLeft()

					}
					if e.Msg.Right {
						p.MoveRight()
					}
				}

				s.mu.Unlock()
			}
		case <-ticker.C:
			s.update()
			s.tick++
			s.broadcast()
		}
	}
}

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
	leftPlayerMsg := PlayerLeaveMsg{
		Type: "playerleft",
		ID:   id,
	}
	data, _ := json.Marshal(leftPlayerMsg)
	for _, p := range s.players {
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
	for _, p := range s.players {
		if p.ID != id && p.Conn != nil {
			p.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}

func (s *Server) broadcastPlayerList() {
	var playerIds []string

	for _, p := range s.players {
		playerIds = append(playerIds, p.ID)
	}

	for _, p := range s.players {
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

	for id, p := range s.players {
		snapshot.Players[id] = PlayerSnapshot{
			Pos: p.Pos,
		}
	}

	data, _ := json.Marshal(snapshot)
	for _, p := range s.players {
		if p.Conn != nil {
			_ = p.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}

func (s *Server) HandleWs(w http.ResponseWriter, r *http.Request) {
	// Check for required parameters BEFORE upgrading connection
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("Connection attempted without a user id. Not upgrading connection")
		http.Error(w, "Missing id parameter", 400)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	p := &player.Player{
		ID:   id,
		Pos:  math.Vec3{X: 0, Y: 0, Z: 0},
		Conn: conn,
	}

	s.msgChannel <- JoinMsg{Type: "join", ID: id, Player: *p}

	go s.playerWsReadLoop(p)
}

func (s *Server) playerWsReadLoop(p *player.Player) {
	defer func() {
		s.msgChannel <- PlayerLeaveMsg{Type: "leave", ID: p.ID}
		if p.Conn != nil {
			_ = p.Conn.Close()
		}
	}()

	for {
		_, data, err := p.Conn.ReadMessage()
		if err != nil {
			return
		}

		var msg WsMsg
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		s.msgChannel <- WsMsg{
			ID:  msg.ID,
			Msg: msg.Msg,
		}
	}
}

func (s *Server) GetPlayer(id string) (*player.Player, bool) {
	player, exists := s.players[id]

	return player, exists
}
