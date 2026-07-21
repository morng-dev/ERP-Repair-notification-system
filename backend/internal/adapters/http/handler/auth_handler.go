package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
	"github.com/morng-dev/erp/pkg/utils"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req entities.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Error:   err.Error(),
		})
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "รูปแบบข้อมูลไม่ถูกต้อง",
			Error:   err.Error(),
		})
	}
	user, err := h.authService.Register(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "ไม่สามารถลงทะเบียนได้",
		})
	}
	return c.Status(fiber.StatusOK).JSON(entities.ApiResponse{
		Success: true,
		Message: "ลงทะเบียนสำเร็จ",
		Data:    user,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req entities.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Error:   err.Error(),
		})
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "รูปแบบข้อมูลไม่ถูกต้อง",
			Error:   err.Error(),
		})
	}
	response, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "เข้าสู่ระบบไม่สำเร็จ",
			Error:   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(entities.ApiResponse{
		Success: true,
		Message: "เข้าสู่ระบบสำเร็จ",
		Data:    response,
	})
}
