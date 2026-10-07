package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
	"github.com/morng-dev/erp/pkg/utils"
)

type AssetHandler struct {
	assetService services.AssetService
}

func NewAssetsHandler(assetService services.AssetService) *AssetHandler {
	return &AssetHandler{assetService: assetService}
}

func (h *AssetHandler) CreateAsset(c *fiber.Ctx) error {
	var id uuid.UUID
	var req entities.CreateAssetRequest
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
	req.OwnerID = id
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
