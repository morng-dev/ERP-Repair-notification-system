package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/morng-dev/erp/internal/config"
)

func main() {
	cfg := config.LoadCongig()
	config.Setupdatabase(cfg)
	config.SetupRedis(cfg)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})

	app.Listen(":3005")
}
