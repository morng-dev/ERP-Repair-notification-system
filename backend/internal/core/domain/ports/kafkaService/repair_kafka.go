package kafkaservice

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type RepairKafkaService interface {
	DeliverRepair(ctx context.Context, repair *entities.Repair) error
}
