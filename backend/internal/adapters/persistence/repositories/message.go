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
		Conteant:   message.Conteant,
		TimeStamp:  time.Now().Unix(),
	}
	if err := r.db.WithContext(ctx).Create(&msg).Error; err != nil {
		return err
	}
	return nil
}
func (r *MessageRepository) GetByID(ctx context.Context, messageID uuid.UUID) (*entities.Message, error) {
	var msg models.Message

	if err := r.db.WithContext(ctx).Preload("Sender").Preload("Receiver").First(&msg, "id = ?", messageID).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(&msg), nil
}

func (r *MessageRepository) GetByChanel(ctx context.Context, chanalID uuid.UUID) (*entities.Message, error) {
	var msgChanel models.Message
	if err := r.db.WithContext(ctx).Where("ChanalID = ?", chanalID).First(&msgChanel).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(&msgChanel), nil
}

func (r *MessageRepository) Delete(ctx context.Context, messageID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Message{}, "id = ?", messageID).Error
}

func (r *MessageRepository) Update(ctx context.Context, messageID uuid.UUID, req *entities.UpdateMessage) error {
	updates := map[string]interface{}{}
	if req.Content != "" {
		updates["Content"] = req.Content
	}
	return r.db.WithContext(ctx).Model(&models.Message{}).Where("id = ?", messageID).Error
}

func (r *MessageRepository) modelsToEntities(msgModel *models.Message) *entities.Message {
	msg := &entities.Message{
		ID:         msgModel.ID,
		SenderID:   msgModel.SenderID,
		ReceiverID: msgModel.ReceiverID,
		ChanalID:   msgModel.ChanalID,
		Conteant:   msgModel.Conteant,
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
	if msgModel.Receiver.ID != uuid.Nil {
		msg.Receiver = &entities.User{
			ID:        msgModel.ReceiverID,
			Email:     msgModel.Receiver.Email,
			Firsname:  msgModel.Receiver.FirstName,
			Lastname:  msgModel.Receiver.LastName,
			Avatar:    msgModel.Receiver.Avatar,
			Active:    msgModel.Receiver.Active,
			CreatedAt: msgModel.Receiver.CreatedAt,
			UpdatedAt: msgModel.Receiver.UpdatedAt,
		}
	}
	return msg
}
