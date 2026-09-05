package services

import (
	"context"
	"errors"

	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
)

type ProfessionService struct {
	professRepo repositories.ProfessionRepository
}

func NewProfessionsService(professRepo repositories.ProfessionRepository) services.ProfressionService {
	return &ProfessionService{professRepo: professRepo}
}

func (s *ProfessionService) CreateProfession(ctx context.Context, req *entities.Profession) (*entities.Profession, error) {
	_, exist := s.professRepo.GetByName(ctx, req.Name)
	if exist == nil {
		return nil, errors.New("profession already exists")
	}
	profession, err := s.professRepo.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return profession, nil
}
