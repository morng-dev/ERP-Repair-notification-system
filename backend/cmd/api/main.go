package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	// app := fiber.New(fiber.Config{
	// 	ErrorHandler: func(c *fiber.Ctx, err error) error {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 			"error": err.Error(),
	// 		})
	// 	},
	// })
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})

	app.Listen(":3005")
}
