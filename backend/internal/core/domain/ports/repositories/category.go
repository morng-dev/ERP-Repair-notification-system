package repositories

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *entities.Category) error
	GetNameExist(ctx context.Context, name string) (bool, error)
}
