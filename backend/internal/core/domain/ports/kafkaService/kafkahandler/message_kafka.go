package kafkahandler

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type MessageHandler interface {
	DeviliryMessage(ctx context.Context, msg *entities.Message) error
}
