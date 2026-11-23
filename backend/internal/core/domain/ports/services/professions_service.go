package services

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type ProfessionService interface {
	CreateProfession(ctx context.Context, req *entities.CreateDescription) (*entities.Profession, error)
	GetProfessionAll(ctx context.Context) ([]*entities.Profession, error)
}
