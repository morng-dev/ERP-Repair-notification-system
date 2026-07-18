package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SenderID uuid.UUID `json:"sender_id"`
	Sender   User      `gorm:"foreignKey:SenderID;references:ID" json:"sender"`

	ReceiverID uuid.UUID `json:"receiver_id"`
	Receiver   User      `gorm:"foreignKey:ReceiverID;references:ID" json:"receiver"`

	ChanalID uuid.UUID `json:"chanal_id"`
	Chanal   Chanal    `gorm:"foreignKey:ChanalID;references:ID" json:"chanal"`

	Messages  string `json:"messages"`
	TimeStamp int64  `json:"timestamp"`
}

type Chanal struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(100)" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   time.Time `gorm:"autoDeleteTime" json:"deleted_at"`

	Members  []User    `gorm:"many2many:user_chanals;" json:"members,omitempty"`
	Messages []Message `gorm:"foreignKey:ChanalID" json:"messages,omitempty"`
}
