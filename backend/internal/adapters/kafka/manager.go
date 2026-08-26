package kafka

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go"
)

type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type KafkaManager struct {
	addr string
}

func NewKafkaManager(addr string) *KafkaManager {
	return &KafkaManager{addr: addr}
}

func (km *KafkaManager) Entopic(topics []string) error {
	kafkaConn, err := kafka.Dial("tcp", km.addr)
	if err != nil {
		return fmt.Errorf("failed to connect kafka %v", err)
	}
	defer kafkaConn.Close()

	controller, err := kafkaConn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get controller %v", err)
	}
	kafkaControllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, fmt.Sprint(controller.Port)))
	if err != nil {
		return fmt.Errorf("failed to connect Kafka controller %v", err)
	}
	defer kafkaControllerConn.Close()

	topicConfig := []kafka.TopicConfig{}

	for _, topic := range topics {
		topicConfig = append(topicConfig, kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		})
	}
	err = kafkaControllerConn.CreateTopics(topicConfig...)
	if err != nil && !errors.Is(err, kafka.TopicAlreadyExists) {
		return fmt.Errorf("failed to create topics: %w", err)
	}
	return nil
}

func (km *KafkaManager) HealthCheck() error {
	conn, err := kafka.Dial("tcp", km.addr)
	if err != nil {
		return fmt.Errorf("fail to connect kafka : %v", err)
	}
	defer conn.Close()

	_, err = conn.Brokers()
	if err != nil {
		return fmt.Errorf("cannot get brockers")
	}
	return nil
}

func WaitKafkaReady(addr string, timeout time.Duration) error {
	km := NewKafkaManager(addr)
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if km.HealthCheck() == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("Kafka not ready after %v", timeout)
}
