package entities

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Notifications struct {
	userID string
	conn   websocket.Conn
	Send   chan interface{}
	mu     sync.RWMutex
	closed bool
}
