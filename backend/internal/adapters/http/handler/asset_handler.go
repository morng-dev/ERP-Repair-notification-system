package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
	"github.com/morng-dev/erp/pkg/utils"
)

type assetHandler struct {
	assetService services.AssetService
}

func NewAssetsHandler(assetService services.AssetService) *assetHandler {
	return &assetHandler{assetService: assetService}
}

func (h *assetHandler) CreateAsset(c *fiber.Ctx) error {
	var req entities.CreateAssetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "Invalid body request",
			Error:   err.Error(),
		})
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.ErrorResponse{
			Success: false,
			Message: "Invalid struct request",
			Error:   err.Error(),
		})
	}

	if err := h.assetService.CreateAsset(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.ErrorResponse{
			Success: false,
			Message: "failed to create asset",
			Error:   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(entities.ApiResponse{
		Success: true,
		Message: "created success",
	})
}
