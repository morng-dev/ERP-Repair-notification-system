package services

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type ProfressionService interface {
	CreateProfession(ctx context.Context, req *entities.Profession) (*entities.Profession, error)
}
