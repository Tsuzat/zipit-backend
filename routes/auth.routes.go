package routes

import (
	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/Tsuzat/zipit-go-fiber/controller"
	"github.com/Tsuzat/zipit-go-fiber/middleware"
)

func InitAuthRouter() {
	group := config.APP.Group("/api/v1/auth")

	group.Post("/signup", controller.SignUpUser)
	group.Post("/login", controller.LoginUser)
	group.Get("/me", middleware.Authenticate, controller.Me)
	group.Get("/logout", controller.LogOut)
	group.Post("/refresh-access-token", controller.RefreshToken)
	group.Get("/verify", controller.VerifyUserEmail)
}
