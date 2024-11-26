package db

import (
	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
)

func InitRedis() error {
	opts, err := redis.ParseURL(config.REDIS_URL)
	if err != nil {
		log.Error("Redis URL parsing error:", err)
		return err
	}
	config.RDB = redis.NewClient(opts)
	return nil
}
