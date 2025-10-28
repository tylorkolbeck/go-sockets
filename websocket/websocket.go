package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tylorkolbeck/go-sockets/internal/logger"
)

type MessageHandler func(conn *websocket.Conn, messageType int, data []byte) error
type ConnectionHandler func(conn *websocket.Conn, r *http.Request) error
type DisconnectHandler func(conn *websocket.Conn)
type Manager struct {
	upgrader          websocket.Upgrader
	messageHandler    MessageHandler
	connectionHandler ConnectionHandler
	disconnectHandler DisconnectHandler
	logger            *logger.Logger
}

func NewWebsocketManager(msgHandler MessageHandler, connHandler ConnectionHandler, disconnectHandler DisconnectHandler) *Manager {
	log := logger.NewLogger("Websocket Manager")
	log.Info("Initializing")
	return &Manager{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		messageHandler:    msgHandler,
		connectionHandler: connHandler,
		disconnectHandler: disconnectHandler,
		logger:            log,
	}
}

func (m *Manager) HandleConnection(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		m.logger.Error("Connection attempted without a user id. Not upgrading connection")
		http.Error(w, "Missing id parameter", 400)
		return
	}

	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		m.logger.Error("Websocket upgrade failed: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	if err := m.connectionHandler(conn, r); err != nil {
		m.logger.Error("Connection handler error: %v", err)
		conn.Close()
		return
	}

	go m.readLoop(conn)
}

func (m *Manager) readLoop(conn *websocket.Conn) {
	defer func() {
		m.disconnectHandler(conn)
		conn.Close()
	}()

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			m.logger.Error("Read error: %v", err)
			break
		}

		// Pass message to handler
		if err := m.messageHandler(conn, messageType, data); err != nil {
			m.logger.Error("Message handler error: %v", err)
			break
		}
	}
}
