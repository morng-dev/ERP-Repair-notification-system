package services

import (
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
)

type MessageService struct {
	messageRepo repositories.MessageRepository
}

func NewMessageService(messageRepo repositories.MessageRepository) services.MessageService {
	return &MessageService{messageRepo: messageRepo}
}
