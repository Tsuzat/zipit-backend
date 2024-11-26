package config

import (
	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
)

var (
	APP *fiber.App
	DB  *pg.DB
)
