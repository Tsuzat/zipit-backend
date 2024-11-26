package db

import (
	"errors"
	"os"

	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/go-pg/pg/v10"
)

func ConnectDB() error {
	dburl := os.Getenv("DB_URL")
	if dburl == "" {
		return errors.New("DB_URL is not set")
	}
	opts, err := pg.ParseURL(dburl)
	if err != nil {
		return err
	}
	config.DB = pg.Connect(opts)
	return nil
}
