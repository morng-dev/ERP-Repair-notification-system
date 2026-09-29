package kafkarepo

import (
	"context"
	"log"

	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	kafkaservice "github.com/morng-dev/erp/internal/core/domain/ports/kafkaService"
	"gorm.io/gorm"
)

type notificationsRepository struct {
	db *gorm.DB
}

func NewnotificationsRepository(db *gorm.DB) kafkaservice.NotificationHandler {
	return &notificationsRepository{db: db}
}

func (k *notificationsRepository) DeliverNotification(msg *entities.Notification) {
	notif := models.Notification{
		UserID:  msg.UserID,
		Message: msg.Content,
		Type:    msg.Type,
		Status:  msg.Status,
	}
	if err := k.db.WithContext(context.Background()).Create(&notif).Error; err != nil {
		log.Printf("Failed to save notification: %v", err)
	} else {
		log.Printf("Notification saved for user %s", msg.UserID)
	}
}
