package consumer

import (
	"time"

	kafkaservice "github.com/morng-dev/erp/internal/core/domain/ports/kafkaService"
	"github.com/segmentio/kafka-go"
)

type NotificationConsumer struct {
	reader  *kafka.Reader
	handler kafkaservice.NotificationService
}

func NewNotificationConsumer(kafkaAddr []string, nodeID string, handler kafkaservice.NotificationService) *NotificationConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           kafkaAddr,
		Topic:             "notification-topic",
		GroupID:           "notification-group-" + nodeID,
		MaxBytes:          10e6,
		HeartbeatInterval: time.Second,
		CommitInterval:    time.Second,
	})
	return &NotificationConsumer{
		reader:  reader,
		handler: handler,
	}
}

// func(c *NotificationConsumer)
