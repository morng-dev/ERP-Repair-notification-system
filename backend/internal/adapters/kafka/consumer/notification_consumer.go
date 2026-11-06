package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/morng-dev/erp/internal/adapters/kafka/retry"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/kafkaService/kafkahandler"

	"github.com/segmentio/kafka-go"
)

type NotificationConsumer struct {
	reader  *kafka.Reader
	handler kafkahandler.NotificationHandler
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewNotificationConsumer(kafkaAddr []string, nodeID string, handler kafkahandler.NotificationHandler) *NotificationConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           kafkaAddr,
		Topic:             "notification-topic",
		GroupID:           "notification-group-" + nodeID,
		MaxBytes:          10e6,
		HeartbeatInterval: time.Second,
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
		kafkaMsg, err := c.reader.FetchMessage(c.ctx)
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
		if event.Type != "notification" {
			log.Printf("unknown event type: %s", event.Type)
			return
		}
		var notification entities.Notification
		if err := json.Unmarshal(event.Data, &notification); err != nil {
			log.Printf("Error unmarshalling notification")
			continue
		}
		err = retry.Do(c.ctx, retry.Config{
			MaxRetries:   10,
			InitialDelay: 1 * time.Second,
			MaxDelay:     1 * time.Second,
		}, func() error {
			return c.handler.HandleNotification(
				c.ctx,
				&notification,
			)

		})
		if err != nil {
			log.Printf("failed to deliver notification after retries: %v", err)
			return
		}
		if err := c.reader.CommitMessages(c.ctx, kafkaMsg); err != nil {
			log.Printf("failed to commit kafka message: %v", err)
			return
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
