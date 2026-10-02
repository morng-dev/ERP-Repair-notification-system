package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/adapters/http/handler"
	"github.com/morng-dev/erp/internal/adapters/http/middleware"
	"github.com/morng-dev/erp/internal/adapters/http/routes"
	"github.com/morng-dev/erp/internal/adapters/kafka"
	mail "github.com/morng-dev/erp/internal/adapters/mailer"
	"github.com/morng-dev/erp/internal/adapters/persistence/kafkarepo"
	"github.com/morng-dev/erp/internal/adapters/persistence/redis"
	"github.com/morng-dev/erp/internal/adapters/persistence/repositories"
	"github.com/morng-dev/erp/internal/config"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/kafkaservice"
	"github.com/morng-dev/erp/internal/core/services"
)

func main() {

	cfg := config.LoadCongig()
	//db
	db := config.Setupdatabase(cfg)
	//redis cache
	clientRedis := config.SetupRedis(cfg)
	//SmtpMail
	mailer := mail.NewSmtpMailer(cfg.SMTP_HOST, cfg.SMTP_USERNAME, cfg.SMTP_PASSWORD, cfg.SMTP_FROM, cfg.APPURL, cfg.SMTP_PORT)
	//redis cache
	userRedisRepo := redis.NewUserRedisRepo(clientRedis)
	//repo
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	profesRepo := repositories.NewProfessionRepository(db)
	permissionRepo := repositories.NewPermissionsRepository(db)
	//kafkaRepo
	notifRepo := kafkarepo.NewnotificationsRepository(db)
	notifService := kafkaservice.NewNotificationsKafkaService(notifRepo)
	//queue
	nm, err := kafka.NewNotificationManager("localhost:29092", "node1", notifService)
	if err != nil {
		log.Fatalf("failed to initialize notification manager: %v", err)
	}

	//middle ware
	authMW := middleware.NewAuthMiddleware(cfg.JWTSecret, clientRedis, permissionRepo)
	//serivce
	authrService := services.NewAuthService(userRepo, roleRepo, userRedisRepo, mailer)
	profesService := services.NewProfessionsService(profesRepo)
	permissionService := services.NewPermissionsService(permissionRepo)

	// messageManager, err := kafka.NewMessageManager(
	// 	"localhost:9092",
	// 	"erp-api-1",
	// 	kafkaHandler{},
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	authHandler := handler.NewAuthHandler(authrService)
	profesHandler := handler.NewProfessionsHandler(profesService)
	permissionHandler := handler.NewPermissionHandler(permissionService)

	app := fiber.New(fiber.Config{
		ProxyHeader:             fiber.HeaderXForwardedFor,
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"10.0.0.0/8", "172.16.0.0/12"},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			msg := "เกิดข้อผิดพลาดภายในระบบ"
			var fe *fiber.Error
			if errors.As(err, &fe) {
				code = fe.Code
				msg = fe.Message
			}
			if code == fiber.StatusInternalServerError {

				log.Printf("[ERROR] %s %s: %v", c.Method(), c.Path(), err)
			}
			return c.Status(code).JSON(fiber.Map{

				"error": msg,
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
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "health ok",
		})
	})
	routes := routes.NewRoutes(
		authMW,
		authHandler,
		profesHandler,
		permissionHandler,
	)
	routes.SetUpRoute(app)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigchan
		log.Println("shutting down graacfully...")
		nm.Close()
		app.ShutdownWithTimeout(15 * time.Second)
	}()
	userID := uuid.MustParse("f4347b29-d4e1-42c9-b19d-bf7064c0731d")
	go func() {
		for i := 0; i < 1000; i++ {
			err = nm.PublishNotification(&entities.Notification{
				UserID:  userID,
				Content: "ทดสอบ Kafka",
				Type:    "info",
				Status:  "unread",
			})

			if err != nil {
				log.Printf("publish failed: %v", err)
			} else {
				log.Printf("published message %d", i+1)
			}

			time.Sleep(15 * time.Second)
		}
	}()
	//start server
	log.Printf("Server starting on port %s", cfg.APPPORT)
	log.Fatal(app.Listen(":" + cfg.APPPORT))
}
