package kafkaservice

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/kafkaService/kafkahandler"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
)

type notificationHandler struct {
	notifiRepo repositories.NotificationsRepository
}

func NewNotificationHandler(notifiRepo repositories.NotificationsRepository) kafkahandler.NotificationHandler {
	return &notificationHandler{notifiRepo: notifiRepo}
}

func (h *notificationHandler) HandleNotification(ctx context.Context, notif *entities.Notification) error {
	if notif.ID == uuid.Nil || notif.UserID == uuid.Nil {
		return fmt.Errorf("invalid notification: missing id or user_id")
	}
	if notif.Status == "" {
		notif.Status = "unread"
	}
	return h.notifiRepo.Create(ctx, notif)
}
