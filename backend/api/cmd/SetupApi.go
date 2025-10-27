package main

import (
	"Server/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SetupAPI(app *fiber.App) {
	log.Println("Settig up Api Routes..")

	// setup routes
	routes.SetupAuthRoutes(app)
	routes.SetupUserRoutes(app)
	routes.SetupPostRoutes(app)
	routes.SetupChatRoutes(app)
	routes.SetupNotificationRoutes(app)

	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "api",
		})
	})

	log.Println("API routes configured successfully")
}
