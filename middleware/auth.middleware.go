package middleware

import (
	"strings"

	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/Tsuzat/zipit-go-fiber/db"
	"github.com/Tsuzat/zipit-go-fiber/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Authenticate(c *fiber.Ctx) error {
	validateJWT(c)
	return nil
}

func validateJWT(c *fiber.Ctx) error {
	var access_token string
	// Find the token in cookies
	access_token = c.Cookies("access_token")
	// if Cookies is not found, find the token in headers
	if access_token == "" {
		access_token = strings.Split(string(c.Request().Header.Peek("Authorization")), "Bearer ")[0]
	}
	// if token is not found, return 403
	if access_token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "No Access Token Provided in Request",
		})
	}
	// decode the token
	token, err := jwt.Parse(access_token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.ACCESS_TOKEN_SECRET), nil
	})
	// If there is an error, return 401
	if err != nil && err != jwt.ErrTokenExpired {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "Error while parsing the access token. Please relogin or refresh your access token",
		})
	}
	if err == jwt.ErrTokenExpired {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "Access token has expired. Please refresh your access token",
		})
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(403).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "Error while parsing the access token",
		})
	}
	// Get the user from the database and attach it to the context so that we can use it in the route
	id, email, name := claims["id"].(string), claims["email"].(string), claims["name"].(string)
	user, err := db.GetUserByIdNameAndEmail(id, name, email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "Invalid token, User not found",
		})
	} else if !user.IsVerified {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized access, User not verified",
		})
	}
	// Attach the user to the context
	c.Locals("user", user)
	return c.Next()
}
