package entities

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID         uuid.UUID `json:"id"`
	SenderID   uuid.UUID `json:"sender_id"`
	Sender     *User     `json:"sender"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Receiver   *User     `json:"receiver"`
	ChanalID   uuid.UUID `json:"chanal_id"`
	Chanal     *Chanal   `json:"chanal"`
	Messages   string    `json:"messages"`
	TimeStamp  int64     `json:"timestamp"`
}

type Chanal struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
