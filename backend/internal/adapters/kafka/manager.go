package kafka

import (
	"encoding/json"
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

func (km *KafkaManager) EnsureTopics(topics []string) error {
	kafkaConn, err := kafka.Dial("tcp", km.addr)
	if err != nil {
		return fmt.Errorf("fail to connected kafka : %v", err)
	}
	defer kafkaConn.Close()
	controller, err := kafkaConn.Controller()
	if err != nil {
		return fmt.Errorf("fail to get controller")
	}
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, fmt.Sprint(controller.Port)))
	if err != nil {
		return fmt.Errorf("fail to connect kafka controller")
	}
	defer controllerConn.Close()
	topicConfig := []kafka.TopicConfig{}

	for _, topic := range topics {
		topicConfig = append(topicConfig, kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		})
	}
	return controllerConn.CreateTopics(topicConfig...)
}

func (km *KafkaManager) HealCheck() error {
	kafkaConn, err := kafka.Dial("tcp", km.addr)
	if err != nil {
		return err
	}
	defer kafkaConn.Close()

	_, err = kafkaConn.Brokers()
	if err != nil {
		return err
	}
	return nil
}

func WaitKafka(addr string, timeout time.Duration) error {
	km := NewKafkaManager(addr)

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if err := km.HealCheck(); err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("Kafka not ready after %v", timeout)
}
