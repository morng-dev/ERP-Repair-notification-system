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

type MessageConsumer struct {
	reader  *kafka.Reader
	handler kafkahandler.MessageHandler
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewMessageConsumer(kafkaAddr []string, nodeID string, handler kafkahandler.MessageHandler) *MessageConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  kafkaAddr,
		GroupID:  "message-group-" + nodeID,
		Topic:    "message-topic",
		MaxBytes: 10e6,
	})
	mc := &MessageConsumer{
		reader:  reader,
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
	}
	return mc
}

func (c *MessageConsumer) Listen() {
	log.Println("kafka message stating")
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}
		KafkaMsg, err := c.reader.FetchMessage(c.ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Printf("error to read message Kafkamessage")
			time.Sleep(time.Second)
			continue
		}
		var event entities.Event
		if err := json.Unmarshal(KafkaMsg.Value, &event); err != nil {
			log.Printf("error to unmarshal message Event consumer")
			continue
		}
		if event.Type != "message" {
			log.Printf("unknown event type: %s", event.Type)
			continue
		}
		var message entities.Message
		if err := json.Unmarshal(event.Data, &message); err != nil {
			log.Printf("error to unmarshal message :%v", err)
			continue
		}
		err = retry.Do(c.ctx, retry.Config{
			MaxRetries:   5,
			InitialDelay: time.Second,
			MaxDelay:     time.Second,
		}, func() error {
			return c.handler.DeviliryMessage(c.ctx, &message)
		})
		if err != nil {
			log.Printf("failed to deliver message after retries: %v", err)
			continue
		}
		if err := c.reader.CommitMessages(c.ctx, KafkaMsg); err != nil {
			log.Printf("failed to commit message: %v", err)
		}
	}
}

func (c *MessageConsumer) Close() {
	c.cancel()
	if c.reader != nil {
		c.reader.Close()
	}
}
