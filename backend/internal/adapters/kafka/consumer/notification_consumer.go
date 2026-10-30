package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/morng-dev/erp/internal/core/domain/entities"
	kafkaservice "github.com/morng-dev/erp/internal/core/domain/ports/kafkaService"
	"github.com/segmentio/kafka-go"
)

type NotificationConsumer struct {
	reader  *kafka.Reader
	handler kafkaservice.NotificationService
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewNotificationConsumer(kafkaAddr []string, nodeID string, handler kafkaservice.NotificationService) *NotificationConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           kafkaAddr,
		Topic:             "notification-topic",
		GroupID:           "notification-group-" + nodeID,
		MaxBytes:          10e6,
		HeartbeatInterval: time.Second,
		CommitInterval:    time.Second,
	})
	nc := &NotificationConsumer{
		reader:  reader,
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
	}
	return nc
}

func (c *NotificationConsumer) Listen() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}
		kafkaMsg, err := c.reader.ReadMessage(c.ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Printf("error to read kafka message :%v", err)
			time.Sleep(time.Second)
			continue
		}
		var event entities.Event
		if err := json.Unmarshal(kafkaMsg.Value, &event); err != nil {
			log.Printf("Error unmarshaling event : %v", err)
			continue
		}
		if event.Type == "notification" {
			var notification entities.Notification
			if err := json.Unmarshal(event.Data, &notification); err != nil {
				log.Printf("Error unmarshalling notification")
				continue
			}
			if err := c.handler.DeliverNotification(c.ctx, &notification); err != nil {
				log.Printf("failed to deliver notification: %v", err)
			}
		}
	}
}

func (c *NotificationConsumer) Close() {
	log.Println("NotificationConsumer stoping")
	c.cancel()
	if c.reader != nil {
		c.reader.Close()
	}
	log.Println("NotificationConsumer stoped")
}
