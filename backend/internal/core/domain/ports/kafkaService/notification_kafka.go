package kafkaservice

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type NotificationService interface {
	DeliverNotification(ctx context.Context, notif *entities.Notification) error
}
