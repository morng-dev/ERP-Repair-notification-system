package services

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type RepairService interface {
	CreateRepair(ctx context.Context, req *entities.RepairRequest) error
}
