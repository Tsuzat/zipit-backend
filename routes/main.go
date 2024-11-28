package routes

import (
	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/Tsuzat/zipit-go-fiber/models"
	"github.com/gofiber/fiber/v2"
)

func checkHealth(c *fiber.Ctx) error {
	return c.Status(200).JSON(models.ApiResponse{
		Status:  200,
		Message: "Status Ok",
		Data:    "Ok",
	})
}

func RountesInit() {
	config.APP.Get("/api/v1/healthcheck", checkHealth)

	// Register Auth Routes
	InitAuthRouter()
	// Register Url Routes
	InitUrlRouter()
}
