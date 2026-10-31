package kafkaproducer

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type NotificationKafkaService interface {
	Publish(ctx context.Context, notif *entities.Notification) error
}
