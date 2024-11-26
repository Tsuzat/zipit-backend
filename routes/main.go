package routes

import (
	"context"

	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/gofiber/fiber/v2"
)

func checkHealth(c *fiber.Ctx) error {
	// Check if the database is connected
	var isDatabase bool
	if config.DB.Ping(context.Background()) != nil {
		isDatabase = true
	} else {
		isDatabase = false
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "ok",
		"database": isDatabase,
	})
}

func RountesInit() {
	config.APP.Get("/api/v1/healthcheck", checkHealth)

	// Register Auth Routes
	InitAuthRouter()
}
