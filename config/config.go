package config

import (
	"os"
	"strconv"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// Environment variables
var PORT string
var DB_URL string
var ACCESS_TOKEN_SECRET string
var ACCESS_TOKEN_EXPIRY int64
var REFRESH_TOKEN_SECRET string
var REFRESH_TOKEN_EXPIRY int64
var BACKEND_URL string

// Email Related
var (
	EMAIL_HOST     string
	EMAIL_PORT     int
	EMAIL_USERNAME string
	EMAIL_PASSWORD string
	EMAIL_PROTOCAL string
	EMAIL_FROM     string
)

// Global variables
var (
	DB  *pg.DB
	APP *fiber.App
)

// Redis Related
var (
	REDIS_URL        string
	REDIS_KEY_EXPIRY int
	RDB              *redis.Client
)

/*
Init function initializes the environment variables
*/
func Init() {
	PORT = os.Getenv("PORT")
	DB_URL = os.Getenv("DB_URL")
	ACCESS_TOKEN_SECRET = os.Getenv("ACCESS_TOKEN_SECRET")
	ACCESS_TOKEN_EXPIRY, _ = strconv.ParseInt(os.Getenv("ACCESS_TOKEN_EXPIRY"), 10, 64)
	REFRESH_TOKEN_SECRET = os.Getenv("REFRESH_TOKEN_SECRET")
	REFRESH_TOKEN_EXPIRY, _ = strconv.ParseInt(os.Getenv("REFRESH_TOKEN_EXPIRY"), 10, 64)
	BACKEND_URL = os.Getenv("BACKEND_URL")

	// Email Related
	EMAIL_HOST = os.Getenv("EMAIL_HOST")
	EMAIL_PORT, _ = strconv.Atoi(os.Getenv("EMAIL_PORT"))
	EMAIL_USERNAME = os.Getenv("EMAIL_USERNAME")
	EMAIL_PASSWORD = os.Getenv("EMAIL_PASSWORD")
	EMAIL_PROTOCAL = os.Getenv("EMAIL_PROTOCAL")
	EMAIL_FROM = os.Getenv("EMAIL_FROM")

	// Redis Related
	REDIS_URL = os.Getenv("REDIS_URL")
	REDIS_KEY_EXPIRY, _ = strconv.Atoi(os.Getenv("REDIS_KEY_EXPIRY"))

}
