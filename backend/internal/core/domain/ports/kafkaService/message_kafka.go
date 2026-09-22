package kafkaservice

import "github.com/morng-dev/erp/internal/core/domain/entities"

type MessageHandler interface {
	DeliverMessage(msg *entities.MessageKafka)
}
