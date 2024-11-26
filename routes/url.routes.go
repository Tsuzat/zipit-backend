package routes

import (
	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/Tsuzat/zipit-go-fiber/controller"
	"github.com/Tsuzat/zipit-go-fiber/middleware"
)

func InitUrlRouter() {
	group := config.APP.Group("/api/v1/url")

	group.Post("/", middleware.Authenticate, controller.CreateUrl)
	group.Get("/", middleware.Authenticate, controller.GetUrls)
	group.Get("/r/:alias", controller.Redirect)
	group.Patch("/:id", middleware.Authenticate, controller.UpdateUrl)
	group.Delete("/:id", middleware.Authenticate, controller.DeleteUrl)
}
