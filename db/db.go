package db

import (
	"errors"

	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/go-pg/pg/v10"
)

func ConnectDB() error {
	if config.DB_URL == "" {
		return errors.New("DB_URL is not set")
	}
	opts, err := pg.ParseURL(config.DB_URL)
	if err != nil {
		return err
	}
	config.DB = pg.Connect(opts)
	if config.DB == nil {
		return errors.New("Could not connect to the database")
	}
	return nil
}
