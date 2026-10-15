package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func (rpm *RepairManager) SendRetries(repair *entities.Repair, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		err := rpm.PublishRepair(repair)
		if err == nil {
			return nil
		}
		log.Printf("Retry %d faild for Repair: %v", i+1, err)
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}
	return fmt.Errorf("faild to sent Repair after %d retry", maxRetries)
}

func (rpm *RepairManager) ListenRepair() {
	log.Println("RepairManager started")
	for {
		select {
		case <-rpm.ctx.Done():
			log.Printf("stop RepairManager")
			return
		default:
		}
		msg, err := rpm.Repairreader.ReadMessage(rpm.ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Printf("error to read Message : %v", err)
			time.Sleep(time.Second)
			continue
		}
		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("error unmarshal event : %v", err)
			continue
		}
		if event.Type == "repair" {
			var repair entities.Repair
			if err := json.Unmarshal(event.Data, &repair); err != nil {
				log.Printf("error unmarshal repair %v", err)
				continue
			}
			if err := rpm.Repaihandler.DeliverRepair(rpm.ctx, &repair); err != nil {
				log.Printf("error DeliverRepair : %v", err)
				continue
			}
		}

	}
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
