package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Order struct {
	EventID   string `json:"event_id"`
	OrderID   string `json:"order_id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

type OrderHandler interface {
	DaliveryOrder(order *Order)
}

type OrderManager struct {
	kafkaWrtier *kafka.Writer
	OrderReader *kafka.Reader
	Handler     OrderHandler
	Ctx         context.Context
	Cancel      context.CancelFunc
}

func NewOrderManager(kafkaAddr string, nodeID string, handler OrderHandler) (*OrderManager, error) {
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
		Topic:          "order",
		GroupID:        "order-group-" + nodeID,
		StartOffset:    kafka.FirstOffset,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	})
	om := &OrderManager{
		kafkaWrtier: writer,
		OrderReader: reader,
		Handler:     handler,
		Ctx:         ctx,
		Cancel:      cancel,
	}
	go om.listenToOrder()
	return om, nil
}

func (om *OrderManager) PublicMessageOrder(msg *Order) error {
	msg.EventID = fmt.Sprintf(msg.OrderID, msg.UserID, msg.Timestamp)

	dataByte, err := json.Marshal(&msg)
	if err != nil {
		return err
	}
	event := &Event{
		Type: "ordermessage",
		Data: json.RawMessage(dataByte),
	}
	eventByte, err := json.Marshal(&event)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	OrderMsg := kafka.Message{
		Topic: "order-message",
		Value: eventByte,
	}
	return om.kafkaWrtier.WriteMessages(ctx, OrderMsg)
}

func (om *OrderManager) listenToOrder() {
	for {
		select {
		case <-om.Ctx.Done():
			return
		default:
		}
		msg, err := om.OrderReader.ReadMessage(om.Ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			log.Printf("kafka Read error: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("error unmarshal event")
			continue
		}
		if event.Type == "order" {
			var orderMsg Order
			if err := json.Unmarshal(event.Data, &orderMsg); err != nil {
				log.Printf("error unmarshal ordermessage")
			}
			om.Handler.DaliveryOrder(&orderMsg)
		}
	}
}
