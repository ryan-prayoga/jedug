package server

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Jedug API v1",
		})
	})
    
    // Placeholder for routes
    // v1.Post("/reports", handlers.CreateReport)
}
