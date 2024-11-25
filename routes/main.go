package routes

import (
	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/gofiber/fiber/v2"
)

func RountesInit() {
	config.APP = fiber.New()

	config.APP.Get("/api/v1/healthcheck", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"status": "ok"})
	})
}
