package repositories

import (
	"context"

	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repositories.CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, category *entities.Category) error {
	categoryModel := &models.Category{
		ID:          category.ID,
		Description: category.Description,
		Image:       category.Image,
	}
	return r.db.WithContext(ctx).Create(&categoryModel).Error
}

func (r *CategoryRepository) GetNameExist(ctx context.Context, name string) (bool, error) {
	var exist bool
	if err := r.db.WithContext(ctx).Raw("SELECT EXISTS(SELECT 1) FROM category WHERE name = ?", name).Scan(&exist).Error; err != nil {
		return false, err
	}
	return exist, nil

}
