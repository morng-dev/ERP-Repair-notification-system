package kafkahandler

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type NotificationConsumerHandler interface {
	HandleNotification(ctx context.Context, notif *entities.Notification) error
}
