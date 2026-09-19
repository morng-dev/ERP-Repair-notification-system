package kafkaservice

import "github.com/morng-dev/erp/internal/core/domain/entities"

type NotificationHandler interface {
	DeliverNotification(notif *entities.Notification)
}
