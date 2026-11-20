package realtime

import (
	"errors"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

var (
	ErrUserOffline    = errors.New("user is offline")
	ErrSendBufferFull = errors.New("notification send buffer is full")
)

type NotificationHub struct {
	mu      sync.RWMutex
	clients map[string]*NotificationClient
	nodeID  string
}

func NewNotificationHub(nodeID string) *NotificationHub {
	return &NotificationHub{
		clients: make(map[string]*NotificationClient),
		nodeID:  nodeID,
	}
}

func (h *NotificationHub) Register(userID uuid.UUID, conn *websocket.Conn) *NotificationClient {

	h.mu.Lock()
	defer h.mu.Unlock()

	userKey := userID.String()

	if oldClient, ok := h.clients[userKey]; ok {
		h.closeClient(oldClient)
		oldClient.Conn.Close()
	}

	client := &NotificationClient{
		UserID: userKey,
		Conn:   conn,
		Send:   make(chan interface{}, 256),
		closed: false,
	}

	h.clients[userKey] = client

	log.Printf("Notification user %s connected on node %s", userKey, h.nodeID)

	return client
}

func (h *NotificationHub) UnregisterClient(userID string, conn *websocket.Conn) *NotificationClient {

	h.mu.Lock()
	defer h.mu.Unlock()

	client, ok := h.clients[userID]
	if !ok {
		return nil
	}

	if client.Conn != conn {
		return nil
	}

	h.closeClient(client)

	delete(h.clients, userID)

	log.Printf("Notification user %s disconnected from node %s", userID, h.nodeID)

	return client
}

func (h *NotificationHub) SendToUser(userID uuid.UUID, message interface{}) error {
	h.mu.RLock()

	client, ok := h.clients[userID.String()]
	if !ok {
		h.mu.RUnlock()
		return ErrUserOffline
	}

	client.mu.Lock()

	if client.closed {
		client.mu.Unlock()
		h.mu.RUnlock()
		return ErrUserOffline
	}

	// ไม่ Block หาก Channel เต็ม
	select {
	case client.Send <- message:
		client.mu.Unlock()
		h.mu.RUnlock()
		return nil

	default:
		client.mu.Unlock()
		h.mu.RUnlock()
		return ErrSendBufferFull
	}
}

func (h *NotificationHub) closeClient(client *NotificationClient) {

	client.mu.Lock()
	defer client.mu.Unlock()

	if client.closed {
		return
	}

	close(client.Send)
	client.closed = true
}
