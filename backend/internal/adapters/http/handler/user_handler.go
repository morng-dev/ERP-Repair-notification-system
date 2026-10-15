package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 0 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	user, Pagination, err := h.userService.GetUserAll(c.Context(), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.ErrorResponse{
			Success: false,
			Message: "can not read user error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(entities.ApiResponse{
		Success:    true,
		Message:    "read data success",
		Data:       user,
		Pagination: Pagination,
	})
}
