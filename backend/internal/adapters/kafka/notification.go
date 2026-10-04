package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	kafkaservice "github.com/morng-dev/erp/internal/core/domain/ports/kafkaService"
	"github.com/segmentio/kafka-go"
)

type NotificationManager struct {
	kafkaWriter  *kafka.Writer
	NotifReader  *kafka.Reader
	NotifHandler kafkaservice.NotificationHandler
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewNotificationManager(kafkaAddr string, nodeID string, handler kafkaservice.NotificationHandler) (*NotificationManager, error) {
	ctx, cancel := context.WithCancel(context.Background())

	writer := &kafka.Writer{
		Addr:         kafka.TCP(kafkaAddr),
		Balancer:     &kafka.Hash{},
		BatchTimeout: 10 * time.Millisecond,
		WriteTimeout: 10 * time.Second,
		RequiredAcks: kafka.RequireOne,
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaAddr},
		Topic:          "notifications-topic",
		GroupID:        "notification-group-" + nodeID,
		MaxBytes:       10e6,
		StartOffset:    kafka.LastOffset,
		CommitInterval: time.Second,
	})
	nm := &NotificationManager{
		kafkaWriter:  writer,
		NotifReader:  reader,
		NotifHandler: handler,
		ctx:          ctx,
		cancel:       cancel,
	}
	log.Println("Notification Manager Initialized")
	go nm.ListenToNotification()
	return nm, nil
}

func (nm *NotificationManager) PublishNotification(notif *entities.Notification) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	notif.CreatedAt = time.Now().UTC()
	notif.ID = uuid.New()

	dataByte, err := json.Marshal(notif)
	if err != nil {
		return err
	}
	event := &Event{
		Type: "notifications",
		Data: dataByte,
	}
	eventByte, err := json.Marshal(event)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s-%s", notif.UserID.String(), notif.Type)
	msg := kafka.Message{
		Topic: "notifications-topic",
		Key:   []byte(key),
		Value: eventByte,
	}
	return nm.kafkaWriter.WriteMessages(ctx, msg)
}

func (nm *NotificationManager) ListenToNotification() {
	log.Println("kafka Notification stating")
	for {
		select {
		case <-nm.ctx.Done():
			return
		default:
		}
		msg, err := nm.NotifReader.ReadMessage(nm.ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Printf("kafka Read error: %v", err)
			time.Sleep(time.Second)
			continue
		}
		log.Printf("topic %v", msg.Topic)
		log.Printf("offse %v", msg.Offset)
		log.Printf("Partition %v", msg.Partition)
		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Error unmarshaling event : %v", err)
			continue
		}
		if event.Type == "notifications" {
			var notif entities.Notification
			if err := json.Unmarshal(event.Data, &notif); err != nil {
				log.Printf("Error unmarshaling notifications : %v", err)
				continue
			}
			if err := nm.NotifHandler.DeliverNotification(nm.ctx, &notif); err != nil {
				log.Printf("failed to deliver notification: %v", err)
				continue
			}
		}

	}
}

func (nm *NotificationManager) Close() {
	log.Println("Stopping notification manager..")
	nm.cancel()

	if nm.kafkaWriter != nil {
		nm.kafkaWriter.Close()
	}
	if nm.NotifReader != nil {
		nm.NotifReader.Close()
	}
	log.Println("notification manager stopped!")
}
