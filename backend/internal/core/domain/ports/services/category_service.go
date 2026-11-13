package services

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req *entities.CreateCategory) error
	GetAllCategory(ctx context.Context, limit, page int) ([]*entities.Category, *entities.PaginationResponse, error)
}
