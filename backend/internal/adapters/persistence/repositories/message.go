package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) repositories.MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, message *entities.Message) error {
	msg := &models.Message{
		SenderID:   message.SenderID,
		ReceiverID: message.ReceiverID,
		ChanalID:   message.ChanalID,
		Messages:   message.Messages,
		TimeStamp:  time.Now().Unix(),
	}
	if err := r.db.WithContext(ctx).Create(&msg).Error; err != nil {
		return err
	}
	return nil
}

func (r *MessageRepository) modelsToEntities(msgModel *models.Message) *entities.Message {
	msg := &entities.Message{
		ID:         msgModel.ID,
		SenderID:   msgModel.SenderID,
		ReceiverID: msgModel.ReceiverID,
		ChanalID:   msgModel.ChanalID,
		Messages:   msgModel.Messages,
		TimeStamp:  msgModel.TimeStamp,
	}
	if msgModel.Sender.ID != uuid.Nil {
		msg.Sender = &entities.User{
			ID:        msgModel.SenderID,
			Email:     msgModel.Sender.Email,
			Firsname:  msgModel.Sender.FirstName,
			Lastname:  msgModel.Sender.LastName,
			Avatar:    msgModel.Sender.Avatar,
			Active:    msgModel.Sender.Active,
			CreatedAt: msgModel.Sender.CreatedAt,
			UpdatedAt: msgModel.Sender.UpdatedAt,
		}
	}
	return msg
}
