package realtime

import (
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type NotificationHub struct {
	mu      sync.RWMutex
	clients map[string]*NotificationClient
	nodeID  string
}

func NewNotificationHub(kafkaAddr []string, nodeID string) *NotificationHub {
	hub := &NotificationHub{
		clients: make(map[string]*NotificationClient),
		nodeID:  nodeID,
	}
	return hub
}

func (h *NotificationHub) Register(userID uuid.UUID, conn *websocket.Conn) *NotificationClient {
	h.mu.Lock()
	defer h.mu.Unlock()

	client := &NotificationClient{
		UserID: userID.String(),
		Conn:   conn,
		Send:   make(chan interface{}, 256),
		closed: false,
	}

	h.clients[userID.String()] = client
	log.Printf("Notification User %s Connected on node %s", userID, h.nodeID)
	return client
}
