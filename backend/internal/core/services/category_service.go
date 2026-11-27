package services

import (
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
)

type CategoryService struct {
	categoryRepo repositories.CategoryRepository
	professRepo  repositories.ProfessionRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository, professRepo repositories.ProfessionRepository) services.CategoryService {
	return &CategoryService{categoryRepo: categoryRepo, professRepo: professRepo}
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *entities.CreateCategory) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	exists, err := s.categoryRepo.GetNameExist(ctx, req.Name)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("category already exists")
	}
	exists, err = s.professRepo.GetByExists(ctx, *req.ProfessionID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("profess notfund")
	}
	category := &entities.Category{
		ID:           id,
		Name:         req.Name,
		Description:  req.Description,
		ProfessionID: req.ProfessionID,
		Image:        req.Image,
	}
	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return errors.New("can not create category")
	}
	return nil
}

func (s *CategoryService) GetAllCategory(ctx context.Context, page, limit int) ([]*entities.Category, *entities.PaginationResponse, error) {
	category, total, err := s.categoryRepo.GetAllCategory(ctx, page, limit)
	if err != nil {
		return nil, nil, err
	}
	totalPage := int(math.Ceil(float64(total) / float64(limit)))
	pagination := &entities.PaginationResponse{
		Page:       page,
		Limit:      limit,
		TotalPages: totalPage,
		TotalItems: total,
	}
	return category, pagination, nil
}
