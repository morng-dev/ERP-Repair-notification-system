package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/morng-dev/erp/internal/adapters/http/handler"
	"github.com/morng-dev/erp/internal/adapters/http/middleware"
)

type Routes struct {
	authMW      *middleware.AuthMiddleware
	authHandler *handler.AuthHandler
}

func NewRoutes(
	authMW *middleware.AuthMiddleware,
	authHandler *handler.AuthHandler,
) *Routes {
	return &Routes{
		authHandler: authHandler,
		authMW:      authMW,
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

	api := app.Group("/api/v1")
	//auth
	auth := api.Group("/auth")
	auth.Post("/register", r.authHandler.Register)
	auth.Post("/login", r.authHandler.Login)

	auth.Get("/", r.authMW.AuthRequire(), r.authHandler.Helloworld)
}
