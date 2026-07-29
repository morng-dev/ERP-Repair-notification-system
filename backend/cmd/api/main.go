package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/morng-dev/erp/internal/adapters/http/handler"
	"github.com/morng-dev/erp/internal/adapters/http/middleware"
	"github.com/morng-dev/erp/internal/adapters/http/routes"
	"github.com/morng-dev/erp/internal/adapters/persistence/repositories"
	"github.com/morng-dev/erp/internal/config"
	"github.com/morng-dev/erp/internal/core/services"
)

func main() {
	cfg := config.LoadCongig()
	db := config.Setupdatabase(cfg)
	redis := config.SetupRedis(cfg)
	authMW := middleware.NewNewAuthMiddleware(cfg.JWTSecret, redis)
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	authrService := services.NewAuthService(userRepo, roleRepo)

	authHandler := handler.NewAuthHandler(authrService)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	routes := routes.NewRoutes(
		authMW,
		authHandler,
	)
	routes.SetUpRoute(app)

	log.Printf("Server starting on port %s", cfg.APPPORT)
	log.Fatal(app.Listen(":" + cfg.APPPORT))
}
