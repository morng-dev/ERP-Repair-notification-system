package kafkaservice

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type RepairKafkaService interface {
	DeliverNotification(ctx context.Context, notif *entities.Repair) error
}
