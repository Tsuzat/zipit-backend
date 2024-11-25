package main

import (
	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/Tsuzat/zipit-go-fiber/routes"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Error("Error loading .env file")
	}

	// Initialize the routes
	routes.RountesInit()

	config.APP.Listen(":8080")
}
