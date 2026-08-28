package kafka

import (
	"context"
	"encoding/json"
	"fmt"
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

type MessageCache struct {
	messages map[string]int64
	mu       sync.RWMutex
}

func NewMessageCache() *MessageCache {
	cache := &MessageCache{
		messages: make(map[string]int64),
	}
	go cache.cleanup()
	return cache
}

func (mc *MessageCache) cleanup() {
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
func (mc *MessageCache) Add(messageID string) bool {
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
	MessageCaChe  *MessageCache
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
	go mm.listenToMessage()
	return mm, nil
}

func (mm *MessageManager) PublicMessage(msg *Message) error {
	msg.MessageID = fmt.Sprintf(msg.FromUserId, msg.ToUSerId, msg.Timestamp)

	dataByte, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	event := Event{
		Type: "message",
		Data: json.RawMessage(dataByte),
	}

	eventByte, err := json.Marshal(event)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	keyConv := GetConv(msg.FromUserId, msg.ToUSerId)
	kafkaMsg := kafka.Message{
		Topic: "message",
		Key:   []byte(keyConv),
		Value: eventByte,
	}
	return mm.kafkaWrtier.WriteMessages(ctx, kafkaMsg)
}

func (mm *MessageManager) listenToMessage() {
	log.Println("kafka message stating")
	for {
		select {
		case <-mm.ctx.Done():
			return
		default:
		}
		msg, err := mm.messageReader.ReadMessage(mm.ctx)
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
			log.Printf("Error unmarshaling event : %v", err)
			continue
		}
		if event.Type == "message" {
			var chatMsg Message
			if err := json.Unmarshal(event.Data, &chatMsg); err != nil {
				log.Printf("Error unmarshaling message : %v", err)
				continue
			}
			if !mm.MessageCaChe.Add(chatMsg.MessageID) {
				log.Printf("Duplicated message Ignored:%s", chatMsg.MessageID)
				continue
			}
			log.Printf("pocess new message from: %s to %s", chatMsg.FromUserId, chatMsg.ToUSerId)
			mm.handler.DeliverMessage(&chatMsg)
		}
	}
}

func GetConv(userID1, userID2 string) string {
	if userID1 > userID2 {
		return userID2 + "-" + userID1
	}
	return userID1 + "-" + userID2
}
