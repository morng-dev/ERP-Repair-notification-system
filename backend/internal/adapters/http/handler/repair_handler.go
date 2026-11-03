package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
	"github.com/morng-dev/erp/pkg/utils"
)

type RepairHandler struct {
	RepairService services.RepairService
}

func NewRepairHandler(RepairService services.RepairService) *RepairHandler {
	return &RepairHandler{RepairService: RepairService}
}

func (h *RepairHandler) CreateRepair(c *fiber.Ctx) error {
	var req entities.RepairRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "invalid request body",
			Error:   err.Error(),
		})
	}
	var id uuid.UUID

	val := c.Locals("userID")
	if val != nil {
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
				Message: "invalid or expiry token",
			})
		}
		Parsed, err := uuid.Parse(claims.UserID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(entities.ErrorResponse{
				Success: false,
				Message: "assetID invalid",
			})
		}
		id = Parsed
	}
	req.ReportedBy = id

	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "Invalid struct request",
			Error:   err.Error(),
		})
	}
	if err := h.RepairService.CreateRepair(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.ErrorResponse{
			Message: "failed to create repair",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(entities.ApiResponse{
		Success: true,
		Message: "created success",
	})
}
