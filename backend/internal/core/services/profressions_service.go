package services

import (
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
)

type ProfessionService struct {
	professRepo repositories.ProfessionRepository
}

func NewProfessionsService(professRepo repositories.ProfessionRepository) services.ProfressionService {
	return &ProfessionService{professRepo: professRepo}
}
