package websocketservice

import (
	"context"
	"errors"

	"github.com/morng-dev/erp/internal/adapters/realtime"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/ws"
)

type NotificationWsService struct {
	notifiHub *realtime.NotificationHub
}

func NewNotificationWsService(notifiHub *realtime.NotificationHub) ws.NotificationWsService {
	return &NotificationWsService{notifiHub: notifiHub}
}

func (ws *NotificationWsService) Deliver(ctx context.Context, notification *entities.Notification) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if notification == nil {
		return errors.New("notification is nil")
	}

	return ws.notifiHub.SendToUser(notification.UserID, notification)
}
