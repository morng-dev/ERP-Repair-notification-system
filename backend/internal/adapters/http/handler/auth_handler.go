package handler

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
			Error:   err.Error(),
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
		return c.Status(fiber.StatusUnauthorized).JSON(entities.ErrorResponse{
			Success: false,
			Message: "เข้าสู่ระบบไม่สำเร็จ",
			Error:   err.Error(),
		})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    response.Token,
		Expires:  time.Now().Add(15 * time.Minute),
		HTTPOnly: true,
		SameSite: "Strict",
	})

	return c.Status(fiber.StatusOK).JSON(entities.ApiResponse{
		Success: true,
		Message: "เข้าสู่ระบบสำเร็จ",
		Data:    response,
	})
}

func (h *AuthHandler) Changepassword(c *fiber.Ctx) error {
	var req entities.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
	}
	var id uuid.UUID
	if val := c.Locals("userID"); val != nil {
		if uid, ok := val.(uuid.UUID); ok {
			id = uid
		}
	}
	if id == uuid.Nil {
		tokenString := c.Cookies("access_token")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(entities.ErrorResponse{
				Success: false,
				Message: "missing access_token cookie",
			})
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(entities.ErrorResponse{
				Success: false,
				Message: "invalid or expired token",
				Error:   err.Error(),
			})
		}

		parsed, err := uuid.Parse(claims.UserID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
				Success: false,
				Message: "id invalid",
				Error:   err.Error(),
			})
		}
		id = parsed
	}
	if err := h.authService.ChangePassword(c.Context(), id, &req); err != nil {
		log.Printf("ChangePassword: service error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "can not change password",
			Error:   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(entities.ApiResponse{
		Success: true,
		Message: "success fully",
	})
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req entities.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
	}
	if err := h.authService.ForgotPassword(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(entities.ApiResponse{
		Success: true,
		Message: "success fully",
	})
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	token := c.Query("token")
	var req entities.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}
	req.ResetToken = token
	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
	}
	if err := h.authService.ResetPassword(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.ErrorResponse{
			Success: false,
			Message: "can't resetPassword",
			Error:   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(entities.ApiResponse{
		Success: true,
		Message: "Reset password success fully",
	})
}
