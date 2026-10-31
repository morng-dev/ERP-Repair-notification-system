package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/kafkaService/kafkaproducer"

	"github.com/segmentio/kafka-go"
)

type NotificationProducer struct {
	writer *kafka.Writer
}

var _ kafkaproducer.NotificationKafkaService = (*NotificationProducer)(nil)

func NewNotificationProducer(Addrkafka ...string) *NotificationProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(Addrkafka...),
		Balancer:     &kafka.Hash{},
		BatchTimeout: 10 * time.Millisecond,
		WriteTimeout: 10 * time.Second,
		RequiredAcks: kafka.RequireAll,
	}
	return &NotificationProducer{writer: writer}
}

func (p *NotificationProducer) Publish(ctx context.Context, notif *entities.Notification) error {
	key := notif.UserID.String()
	dataByte, err := json.Marshal(notif)
	if err != nil {
		return err
	}
	event := &entities.Event{
		Type: "notification",
		Data: dataByte,
	}
	eventByte, err := json.Marshal(event)
	if err != nil {
		return err
	}
	kafkaMsg := kafka.Message{
		Topic: "notification-topic",
		Value: eventByte,
		Key:   []byte(key),
	}
	if err := p.writer.WriteMessages(ctx, kafkaMsg); err != nil {
		return err
	}
	return nil
}

func (p *NotificationProducer) Close() error {
	return p.writer.Close()
}
