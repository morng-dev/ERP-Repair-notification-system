package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
	"github.com/morng-dev/erp/pkg/utils"
)

type ProfessionHandler struct {
	professService services.ProfressionService
	// kafkaManager   *kafka.MessageManager
}

func NewProfessionsHandler(professService services.ProfressionService) *ProfessionHandler {
	return &ProfessionHandler{professService: professService}
}

func (h *ProfessionHandler) CreateProfession(c *fiber.Ctx) error {
	var req entities.Profession
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
	profesions, err := h.professService.CreateProfession(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "Failed to create profession",
			Error:   err.Error(),
		})
	}

	// h.kafkaManager.PublicMessage(&kafka.Message{
	// 	FromUserId: req.Name,
	// 	ToUSerId:   req.Description,
	// })

	return c.Status(fiber.StatusOK).JSON(entities.ApiResponse{
		Success: true,
		Message: "created success fully",
		Data:    profesions,
	})
}
