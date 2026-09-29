package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type NotificationsRepository interface {
	Create(ctx context.Context, req *entities.Notification) error
	UpdateStatus(ctx context.Context, notifID uuid.UUID, req *entities.UpdateStatusNotif) error
	GetByUser(ctx context.Context, userID uuid.UUID) (*entities.Notification, error)
	GetByID(ctx context.Context, notifID uuid.UUID) (*entities.Notification, error)
	Delete(ctx context.Context, notifID uuid.UUID) error
}
