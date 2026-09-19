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
	Conteant   string    `json:"conteant"`
	TimeStamp  int64     `json:"timestamp"`
}

type MessageKafka struct {
	FromUserId string `json:"form_user_id"`
	ToUSerId   string `json:"to_user_id"`
	Content    string `json:"content"`
	Timestamp  string `json:"timestamp"`
	MessageID  string `json:"message_id,omitempty"`
}

type Chanal struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateMessage struct {
	Content string `json:"content"`
}
