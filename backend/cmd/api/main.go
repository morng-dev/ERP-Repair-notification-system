package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/morng-dev/erp/internal/adapters/http/handler"
	"github.com/morng-dev/erp/internal/adapters/http/middleware"
	"github.com/morng-dev/erp/internal/adapters/http/routes"
	"github.com/morng-dev/erp/internal/adapters/persistence/redis"
	"github.com/morng-dev/erp/internal/adapters/persistence/repositories"
	"github.com/morng-dev/erp/internal/config"
	"github.com/morng-dev/erp/internal/core/services"
)

func main() {
	cfg := config.LoadCongig()
	//db
	db := config.Setupdatabase(cfg)
	//redis cache
	clientRedis := config.SetupRedis(cfg)
	//repo
	userRedisRepo := redis.NewUserRedisRepo(clientRedis)
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	profesRepo := repositories.NewProfessionRepository(db)
	permissionRepo := repositories.NewPermissionsRepository(db)
	//middle ware
	authMW := middleware.NewAuthMiddleware(cfg.JWTSecret, clientRedis, permissionRepo)

	authrService := services.NewAuthService(userRepo, roleRepo, userRedisRepo)
	profesService := services.NewProfessionsService(profesRepo)

	authHandler := handler.NewAuthHandler(authrService)
	profesHandler := handler.NewProfessionsHandler(profesService)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	routes := routes.NewRoutes(
		authMW,
		authHandler,
		profesHandler,
	)
	routes.SetUpRoute(app)

	log.Printf("Server starting on port %s", cfg.APPPORT)
	log.Fatal(app.Listen(":" + cfg.APPPORT))
}
