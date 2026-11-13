package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
	"github.com/morng-dev/erp/pkg/utils"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req entities.CreateCategory
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "invalid body request",
		})
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "Invalid struct request",
			Error:   err.Error(),
		})
	}
	if err := h.categoryService.CreateCategory(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.ErrorResponse{
			Success: false,
			Message: "can not create category",
		})
	}
	return c.Status(fiber.StatusOK).JSON(entities.ApiResponse{
		Success: true,
		Message: "create category success",
	})
}
