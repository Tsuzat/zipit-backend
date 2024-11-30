package controller

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/Tsuzat/zipit-go-fiber/db"
	"github.com/Tsuzat/zipit-go-fiber/models"
	"github.com/Tsuzat/zipit-go-fiber/utils"
	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func CreateUrl(c *fiber.Ctx) error {
	var req models.CreateUrlRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ApiError{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body. Please check the request body",
			Error:   err,
		})
	}
	// URL is mandatory
	if strings.Trim(req.Url, " ") == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ApiError{
			Status:  fiber.StatusBadRequest,
			Message: "Url is mandatory",
			Error:   nil,
		})
	}
	// Alias is optional
	if strings.Trim(req.Alias, " ") == "" {
		req.Alias = utils.GenerateRandomString(7)
	}
	foundNewAlias := false
	for i := 0; i < 10 && !foundNewAlias; i++ {
		url, err := db.GetUrlByAlias(req.Alias)
		if err != nil && err != pg.ErrNoRows {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
				Status:  fiber.StatusInternalServerError,
				Message: "An unexpected error occurred. Please try again later.",
				Error:   err,
			})
		}
		if url == nil {
			foundNewAlias = true
			break
		}
		req.Alias = utils.GenerateRandomString(7)
	}
	user := c.Locals("user").(*models.User)
	url := &models.Url{
		Url:       req.Url,
		Alias:     req.Alias,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
		Owner:     user.Id,
	}
	if err := db.InsertUrl(url); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
			Status:  fiber.StatusInternalServerError,
			Message: "An unexpected error occurred. Please try again later.",
			Error:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.ApiResponse{
		Status:  fiber.StatusOK,
		Message: "Url created successfully",
		Data:    url,
	})
}

func GetUrls(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	urls, err := db.GetUrlsByOwnerId(user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
			Status:  fiber.StatusInternalServerError,
			Message: "An unexpected error occurred. Please try again later.",
			Error:   err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(models.ApiResponse{
		Status:  fiber.StatusOK,
		Message: "Urls fetched successfully",
		Data:    urls,
	})
}

func Redirect(c *fiber.Ctx) error {
	// Get the Url Parameter
	alias := c.Params("alias")
	if alias == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ApiError{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body. Please check the request body",
			Error:   nil,
		})
	}
	// Find the alias in redis
	fullUrl, err := config.RDB.Get(context.Background(), alias).Result()
	// If we have a full url in redis, redirect to it
	if err == nil {
		log.Info("Redirecting to full url from redis")
		return c.Redirect(fullUrl)
	}

	// Find the url by alias
	url, err := db.GetUrlByAlias(alias)
	if err != nil && err != pg.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
			Status:  fiber.StatusInternalServerError,
			Message: "An unexpected error occurred. Please try again later.",
			Error:   err,
		})
	}
	// If the url is not found, return 404
	//! TODO: In this case, we can redirect to the frontend 404 page
	if url == nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ApiError{
			Status:  fiber.StatusNotFound,
			Message: "Url not found with alias: " + alias,
			Error:   nil,
		})
	}

	// Save the full url in redis
	err = config.RDB.Set(context.Background(), alias, url.Url,
		time.Second*time.Duration(config.REDIS_KEY_EXPIRY)).Err()

	if err != nil {
		log.Error("Redis set error:", err)
	}

	return c.Redirect(url.Url)
}

func UpdateUrl(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ApiError{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Request URL. Please provide a valid id",
			Error:   nil,
		})
	}

	var req models.UpdateUrlRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ApiError{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err,
		})
	}
	// Check for required fields
	if strings.Trim(req.Url, " ") == "" ||
		strings.Trim(req.Alias, " ") == "" ||
		req.ExpiresAt.IsZero() ||
		req.ExpiresAt.Before(time.Now()) {
		log.Error("Invalid request body", req)
		return c.Status(fiber.StatusBadRequest).JSON(models.ApiError{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Error:   nil,
		})
	}
	// Find the url by id
	url, err := db.GetUrlById(id)
	if err != nil && err != pg.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
			Status:  fiber.StatusInternalServerError,
			Message: "An unexpected error occurred. Please try again later.",
			Error:   err,
		})
	}
	if url == nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ApiError{
			Status:  fiber.StatusNotFound,
			Message: fmt.Sprintf("Url not found with id: %d", id),
			Error:   nil,
		})
	}
	// Update the url
	url.Url = req.Url
	url.Alias = req.Alias
	url.ExpiresAt = req.ExpiresAt
	if err := db.UpdateUrl(url); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
			Status:  fiber.StatusInternalServerError,
			Message: "An unexpected error occurred. Please try again later.",
			Error:   err,
		})
	}
	// Make sure to remove the user from the url
	return c.Status(fiber.StatusOK).JSON(models.ApiResponse{
		Status:  fiber.StatusOK,
		Message: "Url updated successfully",
		Data:    url,
	})
}

func DeleteUrl(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ApiError{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Request URL. Please provide a valid id",
			Error:   nil,
		})
	}

	// Get the url from the database
	url, err := db.GetUrlById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
			Status:  fiber.StatusInternalServerError,
			Message: "Internal server error",
			Error:   err,
		})
	}

	// Make sure the user is the owner of the url
	if url.Owner != c.Locals("user").(*models.User).Id {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   nil,
		})
	}

	// Delete the url from the database
	err = db.DeleteUrl(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
			Status:  fiber.StatusInternalServerError,
			Message: "Internal server error",
			Error:   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.ApiResponse{
		Status:  fiber.StatusOK,
		Message: "Url deleted successfully",
		Data:    nil,
	})
}

func CountUrls(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	count, err := db.CountUrlsByOwnerId(user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
			Status:  fiber.StatusInternalServerError,
			Message: "An unexpected error occurred. Please try again later.",
			Error:   err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(models.ApiResponse{
		Status:  fiber.StatusOK,
		Message: "Urls count fetched successfully",
		Data:    count,
	})
}
