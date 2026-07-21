package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/morng-dev/erp/internal/adapters/http/handler"
	"github.com/morng-dev/erp/internal/adapters/http/routes"
	"github.com/morng-dev/erp/internal/adapters/persistence/repositories"
	"github.com/morng-dev/erp/internal/config"
	"github.com/morng-dev/erp/internal/core/services"
)

func main() {
	cfg := config.LoadCongig()
	db := config.Setupdatabase(cfg)
	config.SetupRedis(cfg)
	userRepo := repositories.NewUserRepository(db)

	authrService := services.NewAuthService(userRepo)

	authHandler := handler.NewAuthHandler(authrService)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	routes := routes.NewRoutes(
		authHandler,
	)
	routes.SetUpRoute(app)

	log.Fatal(app.Listen("server running on port: " + cfg.APPPORT))
}
