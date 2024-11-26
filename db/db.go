package db

import (
	"errors"

	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/Tsuzat/zipit-go-fiber/models"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/gofiber/fiber/v2/log"
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
	// Create the Schemas
	err = createSchema()
	if err != nil {
		log.Error("Error creating database schema: ", err)
		return err
	}
	return nil
}

func createSchema() error {
	models := []interface{}{
		(*models.User)(nil),
		(*models.Url)(nil),
	}
	for _, model := range models {
		err := config.DB.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			log.Error("Error creating table: ", err)
			return err
		}
	}
	return nil
}
