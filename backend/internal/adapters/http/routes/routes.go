package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/morng-dev/erp/internal/adapters/http/handler"
)

type Routes struct {
	authHandler *handler.AuthHandler
}

func NewRoutes(authHandler *handler.AuthHandler) *Routes {
	return &Routes{
		authHandler: authHandler,
	}
}

func (r *Routes) SetUpRoute(app *fiber.App) {
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	auth := app.Group("/api/v1")
	auth.Post("/register", r.authHandler.Register)
	auth.Post("/login", r.authHandler.Login)
}
