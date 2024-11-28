package main

import (
	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/Tsuzat/zipit-go-fiber/db"
	"github.com/Tsuzat/zipit-go-fiber/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Error("Error loading .env file")
	}

	config.Init()

	// Connect to the database
	if err := db.ConnectDB(); err != nil {
		log.Error("Error connecting to the database", err)
		return
	}
	defer config.DB.Close()

	// Connect to Redis
	if err := db.InitRedis(); err != nil {
		log.Error("Error connecting to Redis", err)
		return
	}
	defer config.RDB.Close()

	config.APP = fiber.New()
	if config.APP == nil {
		log.Error("Error creating the app")
		return
	}
	// Initialize the routes
	routes.RountesInit()

	config.APP.Listen(":8080")
}
