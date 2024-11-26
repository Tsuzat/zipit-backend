package db

import (
	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/Tsuzat/zipit-go-fiber/models"
	"github.com/gofiber/fiber/v2/log"
)

func GetUrlByAlias(alias string) (*models.Url, error) {
	url := &models.Url{
		Alias: alias,
	}
	if err := config.DB.Model(url).Where("alias = ?", alias).Select(); err != nil {
		return nil, err
	} else if url.Id == "" {
		return nil, nil
	}
	return url, nil
}

func GetUrlById(id string) (*models.Url, error) {
	url := &models.Url{
		Id: id,
	}
	if err := config.DB.Model(url).WherePK().Select(); err != nil {
		return nil, err
	} else if url.Id == "" {
		return nil, nil
	}
	return url, nil
}

func InsertUrl(url *models.Url) error {
	_, err := config.DB.Model(url).Insert()
	if err != nil || url.Id == "" {
		return err
	}
	log.Info("Created Url with id: ", url.Id)
	return nil
}

func UpdateUrl(url *models.Url) error {
	_, err := config.DB.Model(url).Where("id = ?", url.Id).Update()
	if err != nil || url.Id == "" {
		return err
	}
	log.Info("Updated Url with id: ", url.Id)
	return nil
}

func DeleteUrl(id string) error {
	res, err := config.DB.Model(&models.Url{}).Where("id = ?", id).Delete()
	if err != nil {
		return err
	}
	if res.RowsAffected() == 1 {
		log.Info("Deleted Url with id: ", id)
	}
	return nil
}

func GetUrlsByOwnerId(ownerId string) ([]models.Url, error) {
	var urls []models.Url
	if err := config.DB.Model(&urls).Where("owner_id = ?", ownerId).Select(); err != nil {
		return nil, err
	}
	return urls, nil
}
