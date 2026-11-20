package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/adapters/realtime"
)

type NotificationTestHandler struct {
	notifHub *realtime.NotificationHub
}

func NewNotificationTestHandler(notifHub *realtime.NotificationHub) *NotificationTestHandler {
	return &NotificationTestHandler{
		notifHub: notifHub,
	}
}

func (h *NotificationTestHandler) SendTest(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Query("userID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid user ID",
		})
	}

	err = h.notifHub.SendToUser(userID, fiber.Map{
		"type":    "test",
		"message": "Hello WebSocket!",
	},
	)

	if err != nil {
		log.Println("Send notification error:", err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true, "message": "notification sent",
	})
}
