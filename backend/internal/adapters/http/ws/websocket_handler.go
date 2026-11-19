package ws

import (
	"github.com/morng-dev/erp/internal/adapters/realtime"
)

type WebSocketHandler struct {
	notifHub *realtime.NotificationHub
}

func NewWebSocketHandler(notifHub *realtime.NotificationHub) *WebSocketHandler {
	return &WebSocketHandler{
		notifHub: notifHub,
	}
}
