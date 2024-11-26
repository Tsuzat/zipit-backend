package db

import (
	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/Tsuzat/zipit-go-fiber/models"
	"github.com/gofiber/fiber/v2/log"
)

func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{
		Email: email,
	}
	if err := config.DB.Model(user).Where("email = ?", email).Select(); err != nil {
		return nil, err
	}
	if user.Id == "" {
		return nil, models.UserNotFound
	}
	return user, nil
}

func InsertUser(user *models.User) error {
	_, err := config.DB.Model(user).Insert()
	if err != nil || user.Id == "" {
		return err
	}
	log.Info("Created User with id: ", user.Id)
	return nil
}

func UpdateUser(user *models.User) error {
	_, err := config.DB.Model(user).Where("id = ?", user.Id).Update()
	if err != nil || user.Id == "" {
		return err
	}
	log.Info("Updated User with id: ", user.Id)
	return nil
}

func GetUserByIdNameAndEmail(id, name, email string) (*models.User, error) {
	user := &models.User{
		Id:    id,
		Name:  name,
		Email: email,
	}
	if err := config.DB.Model(user).WherePK().Select(); err != nil {
		return nil, err
	}
	return user, nil
}
