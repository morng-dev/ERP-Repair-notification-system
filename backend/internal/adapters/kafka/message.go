package kafka

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	FromUserId string `json:"form_user_id"`
	ToUSerId   string `json:"to_user_id"`
	Content    string `json:"content"`
	Timestamp  string `json:"timestamp"`
	MessageID  string `json:"message_id,omitempty"`
}

type MessageCaChe struct {
	messages map[string]int64
	mu       sync.RWMutex
}

func NewMessageCache() *MessageCaChe {
	cache := &MessageCaChe{
		messages: make(map[string]int64),
	}
	go cache.cleanup()
	return cache
}

func (mc *MessageCaChe) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		mc.mu.Lock()
		now := time.Now().Unix()
		for id, ts := range mc.messages {
			if now-ts > 300 {
				delete(mc.messages, id)
			}
		}
		mc.mu.Unlock()
	}
}
func (mc *MessageCaChe) Add(messageID string) bool {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	if _, exits := mc.messages[messageID]; exits {
		return false
	}
	mc.messages[messageID] = time.Now().Unix()
	return true
}

type MessageHandler interface {
	DeliverMessage(msg *Message)
}

type MessageManager struct {
	kafkaWrtier   *kafka.Writer
	messageReader *kafka.Reader
	MessageCaChe  *MessageCaChe
	handler       MessageHandler
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewMessageManager(kafaAddr string, nodeID string, handler MessageHandler) (*MessageManager, error) {
	ctx, cancel := context.WithCancel(context.Background())
	writer := &kafka.Writer{
		Addr:         kafka.TCP(kafaAddr),
		Balancer:     &kafka.Hash{},
		BatchTimeout: 10 * time.Millisecond,
		WriteTimeout: 10 * time.Second,
		RequiredAcks: kafka.RequireOne,
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafaAddr},
		Topic:          "test-topic",
		GroupID:        "test-group-" + nodeID,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	})

	mm := &MessageManager{
		kafkaWrtier:   writer,
		messageReader: reader,
		MessageCaChe:  NewMessageCache(),
		handler:       handler,
		ctx:           ctx,
		cancel:        cancel,
	}
	log.Println("Message Manager Initialized")
	return mm, nil
}
