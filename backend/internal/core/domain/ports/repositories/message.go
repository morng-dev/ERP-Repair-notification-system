package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type MessageRepository interface {
	Create(ctx context.Context, message *entities.Message) error
	GetByID(ctx context.Context, messageID uuid.UUID) (*entities.Message, error)
	GetByChanel(ctx context.Context, chanalID uuid.UUID) (*entities.Message, error)
	Update(ctx context.Context, messageID uuid.UUID, req *entities.UpdateMessage) error
	Delete(ctx context.Context, messageID uuid.UUID) error
}
