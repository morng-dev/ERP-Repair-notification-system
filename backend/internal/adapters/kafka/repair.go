package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/morng-dev/erp/internal/core/domain/entities"
	kafkaservice "github.com/morng-dev/erp/internal/core/domain/ports/kafkaService"
	"github.com/segmentio/kafka-go"
)

type RepairManager struct {
	kafkaWriter  *kafka.Writer
	Repairreader *kafka.Reader
	Repaihandler kafkaservice.RepairKafkaService
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewRepairManager(kafkaAddr string, nodeID string, handler kafkaservice.RepairKafkaService) (*RepairManager, error) {
	ctx, cancel := context.WithCancel(context.Background())
	writer := &kafka.Writer{
		Addr:         kafka.TCP(kafkaAddr),
		Balancer:     &kafka.Hash{},
		BatchTimeout: 10 * time.Millisecond,
		WriteTimeout: 10 * time.Second,
		RequiredAcks: kafka.RequireOne,
	}
	reader := kafka.NewReader(kafka.ReaderConfig{})
	rpm := &RepairManager{
		kafkaWriter:  writer,
		Repairreader: reader,
		Repaihandler: handler,
		ctx:          ctx,
		cancel:       cancel,
	}
	return rpm, nil
}

func (rpm *RepairManager) PublishRepair(repair *entities.Repair) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	key := repair.AssetID.String()
	dataByte, err := json.Marshal(repair)
	if err != nil {
		log.Printf("error mashal Repair %v", err)
	}

	event := Event{
		Type: "repair",
		Data: dataByte,
	}

	eventByte, err := json.Marshal(event)
	if err != nil {
		log.Printf("error mashal eventByte %v", err)
	}
	msg := kafka.Message{
		Topic: "repair-topic",
		Key:   []byte(key),
		Value: eventByte,
	}
	return rpm.kafkaWriter.WriteMessages(ctx, msg)
}

func (rpm *RepairManager) Close() {
	log.Printf("Stopping RepairManager")
	rpm.cancel()
	if rpm.kafkaWriter != nil {
		rpm.kafkaWriter.Close()
	}
	if rpm.Repairreader != nil {
		rpm.Repairreader.Close()
	}
	log.Println("repair manager stopped!")
}
