package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/ws"

	"github.com/segmentio/kafka-go"
)

type NotificationWsConsumer struct {
	reader  *kafka.Reader
	handler ws.NotificationWsService
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewNotificationWsConsumer(kafkaAddr []string, nodeID string, handler ws.NotificationWsService) *NotificationWsConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           kafkaAddr,
		Topic:             "notification-topic",
		GroupID:           "notification-group-" + nodeID,
		MaxBytes:          10e6,
		HeartbeatInterval: time.Second,
	})
	nws := &NotificationWsConsumer{
		reader:  reader,
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
	}
	return nws
}

func (c *NotificationWsConsumer) Listen() {
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
			log.Printf("skip unsupported event type: %s", event.Type)

			if err := c.reader.CommitMessages(c.ctx, kafkaMsg); err != nil {
				log.Printf("failed to commit skipped message: %v", err)
				return
			}

			continue
		}
		var notification entities.Notification
		if err := json.Unmarshal(event.Data, &notification); err != nil {
			log.Printf("Error unmarshalling notification")
			continue
		}

		err = c.handler.Deliver(c.ctx, &notification)
		if err != nil {
			log.Printf("failed to deliver notification: %v", err)
		}
		if err := c.reader.CommitMessages(c.ctx, kafkaMsg); err != nil {
			log.Printf("failed to commit kafka message: %v", err)
			return
		}
	}
}

func (c *NotificationWsConsumer) Close() {
	log.Println("NotificationConsumer stoping")
	c.cancel()
	if c.reader != nil {
		if err := c.reader.Close(); err != nil {
			log.Printf("failed to close kafka reader: %v", err)
		}
	}
	log.Println("NotificationConsumer stoped")
}
