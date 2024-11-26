package controller

import (
	"strings"
	"time"

	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/Tsuzat/zipit-go-fiber/db"
	"github.com/Tsuzat/zipit-go-fiber/mail"
	"github.com/Tsuzat/zipit-go-fiber/models"
	"github.com/Tsuzat/zipit-go-fiber/utils"
	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUpUser(c *fiber.Ctx) error {
	var req models.UserSignupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ApiError{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body. Please provide valid input.",
			Error:   err,
		})
	}
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Password hashing error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to secure your password. Please try again.",
			Error:   err,
		})
	}
	req.Password = string(hashedPassword)
	// Check if user already exists
	user, err := db.GetUserByEmail(req.Email)
	if err != nil && err != pg.ErrNoRows {
		log.Warn("Database error when checking user existence:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
			Status:  fiber.StatusInternalServerError,
			Message: "An unexpected error occurred. Please try again later.",
			Error:   err,
		})
	}
	if user != nil && len(user.Id) > 0 && user.IsVerified {
		return c.Status(fiber.StatusConflict).JSON(models.ApiError{
			Status:  fiber.StatusConflict,
			Message: "A verified account already exists with this email. Please log in.",
			Error:   models.UserAlreadyExists,
		})
	} else if user != nil && len(user.Id) > 0 && !user.IsVerified {
		// User already exists but is not verified
		// so we update the user and send verification email
		user.Password = req.Password
		user.Name = req.Name
		user.VerificationToken = utils.GenerateRandomString(50)
		user.VerificationTokenExpiry = time.Now().Add(30 * time.Minute)
		err = db.UpdateUser(user)
		if err != nil {
			log.Error("User update error:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
				Status:  fiber.StatusInternalServerError,
				Message: "An unexpected error occurred. Please try again later.",
				Error:   err,
			})
		}
		go mail.SendEmailVerification(user)
		return c.Status(fiber.StatusConflict).JSON(models.ApiError{
			Status:  fiber.StatusConflict,
			Message: "An unverified account already exists with this email. Please check your email for verification link.",
			Error:   models.UserAlreadyExists,
		})
	}
	// Create the user
	user = &models.User{
		Name:                    req.Name,
		Email:                   req.Email,
		Password:                req.Password,
		VerificationToken:       utils.GenerateRandomString(50),
		VerificationTokenExpiry: time.Now().Add(7 * 24 * time.Hour),
	}
	_, err = config.DB.Model(user).Insert()
	if err != nil {
		log.Error("User creation error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to create your account. Please try again later.",
			Error:   err,
		})
	}
	go mail.SendEmailVerification(user)
	// Successful response
	return c.Status(fiber.StatusCreated).JSON(models.ApiResponse{
		Status:  fiber.StatusCreated,
		Message: "User created successfully. Please verify your email to activate your account.",
		Data:    nil,
	})
}

func LoginUser(c *fiber.Ctx) error {
	var req models.UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ApiError{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err,
		})
	}
	user, err := db.GetUserByEmail(req.Email)
	if err != nil && err != pg.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
			Status:  fiber.StatusInternalServerError,
			Message: "Internal server error while fetching user details",
			Error:   err,
		})
	} else if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "User not found with email",
			Error:   err,
		})
	} else if user.IsVerified == false {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "User not verified with email. Please verify your email or register again.",
			Error:   err,
		})
	} else if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "Invalid password for email",
			Error:   err,
		})
	}

	// Create Access Token
	accessToken, err := user.GenerateJWTToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
			Status:  fiber.StatusInternalServerError,
			Message: "Internal server error while generating access token",
			Error:   err,
		})
	}
	// Create Refresh Token
	refreshToken, err := user.GenerateRefreshToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
			Status:  fiber.StatusInternalServerError,
			Message: "Internal server error while generating refresh token",
			Error:   err,
		})
	}
	// Update User
	user.RefreshToken = refreshToken
	if err := db.UpdateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
			Status:  fiber.StatusInternalServerError,
			Message: "Internal server error while updating user with refresh token",
			Error:   err,
		})
	}
	// Setup Cookies
	c.Cookie(&fiber.Cookie{
		Name:  "access_token",
		Value: accessToken,
	})
	c.Cookie(&fiber.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
	})

	return c.Status(fiber.StatusOK).JSON(models.ApiResponse{
		Status:  fiber.StatusOK,
		Message: "User login successful",
		Data: models.UserLoginResponse{
			Message: "User login successful with id: " + user.Id,
		},
	})
}

func Me(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	return c.Status(fiber.StatusOK).JSON(models.ApiResponse{
		Status:  fiber.StatusOK,
		Message: "User fetched successfully",
		Data:    "Hello There!! " + user.Name,
	})
}

func LogOut(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:  "access_token",
		Value: "",
	})
	c.Cookie(&fiber.Cookie{
		Name:  "refresh_token",
		Value: "",
	})
	return c.Status(fiber.StatusOK).JSON(models.ApiResponse{
		Status:  fiber.StatusOK,
		Message: "Logged out successfully",
		Data:    nil,
	})
}

func RefreshToken(c *fiber.Ctx) error {
	var refresh_token string
	// Find the token in cookies
	refresh_token = c.Cookies("refresh_token")
	// if Cookies is not found, find the token in headers
	if refresh_token == "" {
		refresh_token = strings.Split(string(c.Request().Header.Peek("Authorization")), "Bearer ")[0]
	}
	// if token is not found, return 403
	if refresh_token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "No Refresh Token Provided in Request",
		})
	}
	// decode the token
	token, err := jwt.Parse(refresh_token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.REFRESH_TOKEN_SECRET), nil
	})
	// If there is an error, return 401
	if err != nil && err != jwt.ErrTokenExpired {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "Could not parse the refresh token. Please relogin or refresh your access token",
		})
	} else if err == jwt.ErrTokenExpired {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "Refresh token has expired. Please login again",
		})
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "Could not parse the refresh token. Please relogin or refresh your access token",
		})
	}
	userId := claims["id"].(string)
	user, err := db.GetUserByIdNameAndEmail(userId, "", "")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "Could not find user with provided refresh token",
		})
	}
	// Check if the refresh token is valid
	if user.RefreshToken != refresh_token {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "Could not find refresh token for the user. Please relogin",
		})
	}
	// Generate new access token
	access_token, err := user.GenerateJWTToken()
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "Could not generate new access token",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:  "access_token",
		Value: access_token,
	})

	return c.Status(fiber.StatusOK).JSON(models.ApiResponse{
		Status:  fiber.StatusOK,
		Message: "Access Token Refreshed successfully",
	})
}

func VerifyUserEmail(c *fiber.Ctx) error {
	// Retrieve Query Parameters
	email := c.Query("email")
	token := c.Query("verification_token")
	// Find the user by email
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
			Status:  fiber.StatusInternalServerError,
			Message: "Could not find user with provided email",
			Error:   err,
		})
	}
	if user.IsVerified {
		return c.Status(fiber.StatusOK).JSON(models.ApiResponse{
			Status:  fiber.StatusOK,
			Message: "User is already verified. Please login to continue",
			Data:    "Hello There!! " + user.Name,
		})
	}
	// Verify the verification token
	if user.VerificationToken != token {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "Invalid verification token",
			Error:   nil,
		})
	}
	// Check if the token has expired
	if time.Now().After(user.VerificationTokenExpiry) {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiError{
			Status:  fiber.StatusUnauthorized,
			Message: "Verification token has expired. Please register again",
			Error:   nil,
		})
	}
	// Update the user
	user.IsVerified = true
	if err := db.UpdateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiError{
			Status:  fiber.StatusInternalServerError,
			Message: "Internal server error while updating user verification status",
			Error:   err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(models.ApiResponse{
		Status:  fiber.StatusOK,
		Message: "Verification successful. Please login to continue",
		Data:    "Hello There!! " + user.Name,
	})
}
