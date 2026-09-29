package kafkaservice

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type MessageHandler interface {
	DeliverMessage(ctx context.Context, msg *entities.MessageKafka)
}
