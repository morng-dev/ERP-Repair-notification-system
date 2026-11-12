package kafkaproducer

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type MessageKafkaService interface {
	Publish(ctx context.Context, message *entities.Message) error
}
