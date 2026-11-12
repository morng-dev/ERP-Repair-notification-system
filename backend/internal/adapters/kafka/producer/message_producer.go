package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/kafkaService/kafkaproducer"
	"github.com/segmentio/kafka-go"
)

type MessageProducer struct {
	writer *kafka.Writer
}

var _ kafkaproducer.MessageKafkaService = (*MessageProducer)(nil)

func NewMessageProducer(kafkaAddr []string) *MessageProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(kafkaAddr...),
		Balancer:     &kafka.Hash{},
		WriteTimeout: 10 * time.Second,
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}
	return &MessageProducer{writer: writer}
}

func (p *MessageProducer) Publish(ctx context.Context, message *entities.Message) error {
	dataByte, err := json.Marshal(message)
	if err != nil {
		return err
	}
	event := &entities.Event{
		Type: "message",
		Data: dataByte,
	}
	eventByte, err := json.Marshal(event)
	if err != nil {
		return err
	}
	keyConv := GetConv(message.SenderID.String(), message.ReceiverID.String())
	kafkaMsg := kafka.Message{
		Topic: "message-topic",
		Key:   []byte(keyConv),
		Value: eventByte,
	}
	return p.writer.WriteMessages(ctx, kafkaMsg)
}

func (p *MessageProducer) Close() error {
	return p.writer.Close()
}

func GetConv(userID1, userID2 string) string {
	if userID1 > userID2 {
		return userID2 + "-" + userID1
	}
	return userID1 + "-" + userID2
}
