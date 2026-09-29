package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"gorm.io/gorm"
)

type NotificationsRepository struct {
	db *gorm.DB
}

func NewNotificationsRepository(db *gorm.DB) repositories.NotificationsRepository {
	return &NotificationsRepository{db: db}
}

func (r *NotificationsRepository) Create(ctx context.Context, notification *entities.Notification) error {
	notif := models.Notification{
		UserID:  notification.UserID,
		Message: notification.Content,
		Type:    notification.Type,
		Status:  notification.Status,
	}
	return r.db.WithContext(ctx).Create(&notif).Error
}

func (r *NotificationsRepository) UpdateStatus(ctx context.Context, notifID uuid.UUID, status *entities.UpdateStatusNotif) error {
	return r.db.WithContext(ctx).Model(&models.Notification{}).Where("id = ?", notifID).Update("status = ?", status).Error
}
func (r *NotificationsRepository) GetByUser(ctx context.Context, userID uuid.UUID) (*entities.Notification, error) {
	var notif models.Notification
	if err := r.db.WithContext(ctx).First(notif, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &entities.Notification{
		ID:        notif.ID,
		UserID:    notif.UserID,
		Content:   notif.Message,
		Type:      notif.Type,
		Status:    notif.Status,
		CreatedAt: notif.CreatedAt,
	}, nil
}
func (r *NotificationsRepository) GetByID(ctx context.Context, notifID uuid.UUID) (*entities.Notification, error) {
	var notif models.Notification
	if err := r.db.WithContext(ctx).First(notif, "id = ?", notifID).Error; err != nil {
		return nil, err
	}
	return &entities.Notification{
		ID:        notif.ID,
		UserID:    notif.UserID,
		Content:   notif.Message,
		Type:      notif.Type,
		Status:    notif.Status,
		CreatedAt: notif.CreatedAt,
	}, nil
}
func (r *NotificationsRepository) Delete(ctx context.Context, notifID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Notification{}, "id = ?", notifID).Error
}
