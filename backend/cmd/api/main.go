package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/morng-dev/erp/internal/adapters/persistence/repositories"
	"github.com/morng-dev/erp/internal/config"
	"github.com/morng-dev/erp/internal/core/services"
)

func main() {
	cfg := config.LoadCongig()
	db := config.Setupdatabase(cfg)
	config.SetupRedis(cfg)
	userRepo := repositories.NewUserRepository(db)

	userService := services.NewAuthService(userRepo)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})

	app.Listen(":3005")
}
