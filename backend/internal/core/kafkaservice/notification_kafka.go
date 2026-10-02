package kafkaservice

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
	kafkaservice "github.com/morng-dev/erp/internal/core/domain/ports/kafkaService"
)

type NotificationsService struct {
	notifHanler kafkaservice.NotificationHandler
}

func NewNotificationsKafkaService(notifHanler kafkaservice.NotificationHandler) kafkaservice.NotificationHandler {
	return &NotificationsService{notifHanler: notifHanler}
}

func (nf *NotificationsService) DeliverNotification(ctx context.Context, notif *entities.Notification) error {
	if err := nf.notifHanler.DeliverNotification(ctx, notif); err != nil {
		return err
	}
	return nil

}
