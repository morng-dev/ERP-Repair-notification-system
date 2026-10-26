package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
)

type CategoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) services.CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *entities.CategoryCreate) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	exists, err := s.categoryRepo.GetNameExist(ctx, req.Name)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Role already exists")
	}
	category := &entities.Category{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
	}
	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return errors.New("can not create category")
	}
	return nil
}
