package ws

import (
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/adapters/realtime"
)

const (
	maxMessageSize = 512
	pongWait       = 60 * time.Second
	pingPeriod     = 54 * time.Second
)

type WebSocketHandler struct {
	notifHub *realtime.NotificationHub
}

func NewWebSocketHandler(notifHub *realtime.NotificationHub) *WebSocketHandler {
	return &WebSocketHandler{
		notifHub: notifHub,
	}
}

func (h *WebSocketHandler) Handle(conn *websocket.Conn) {
	userIDValue := conn.Locals("userID")

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		conn.Close()
		return
	}
	client := h.notifHub.Register(userID, conn)
	// Write Pump ทำหน้าที่เขียนข้อมูลเพียง Goroutine เดียว
	go h.writePump(client)
	// Read Pump ทำหน้าที่อ่านข้อมูลและรอ Disconnect
	defer func() {
		h.notifHub.UnregisterClient(userID.String(), conn)
		conn.Close()
	}()

	h.readPump(client)
}

func (h *WebSocketHandler) readPump(client *realtime.NotificationClient) {
	conn := client.Conn

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))

	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func (h *WebSocketHandler) writePump(client *realtime.NotificationClient) {
	conn := client.Conn

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				return
			}

			if err := conn.WriteJSON(message); err != nil {
				return
			}

		case <-ticker.C:
			if err := conn.WriteMessage(
				websocket.PingMessage,
				nil,
			); err != nil {
				return
			}
		}
	}
}
