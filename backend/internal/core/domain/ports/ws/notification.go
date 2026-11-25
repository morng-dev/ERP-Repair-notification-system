package ws

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type NotificationWsHandler interface {
	Deliver(ctx context.Context, notification *entities.Notification) error
}
