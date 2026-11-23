package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
)

type ProfessionService struct {
	professRepo repositories.ProfessionRepository
}

func NewProfessionsService(professRepo repositories.ProfessionRepository) services.ProfessionService {
	return &ProfessionService{professRepo: professRepo}
}

func (s *ProfessionService) CreateProfession(ctx context.Context, req *entities.CreateDescription) (*entities.Profession, error) {
	_, exist := s.professRepo.GetByName(ctx, req.Name)
	if exist == nil {
		return nil, errors.New("profession already exists")
	}
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	professions := &entities.Profession{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}

	profession, err := s.professRepo.Create(ctx, professions)
	if err != nil {
		return nil, err
	}
	return profession, nil
}

func (s *ProfessionService) GetProfessionAll(ctx context.Context) ([]*entities.Profession, error) {
	profession, err := s.professRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return profession, nil
}
