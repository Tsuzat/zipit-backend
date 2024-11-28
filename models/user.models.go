package models

import (
	"time"

	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Id                      int       `json:"id" pg:"id"`
	Name                    string    `json:"name" pg:"name"`
	Email                   string    `json:"email" pg:"email"`
	Password                string    `json:"password" pg:"password"`
	ProfileImage            string    `json:"profile_image" pg:"profile_image"`
	IsVerified              bool      `json:"is_verified" pg:"is_verified"`
	RefreshToken            string    `json:"refresh_token" pg:"refresh_token"`
	VerificationToken       string    `json:"verification_token" pg:"verification_token"`
	VerificationTokenExpiry time.Time `json:"verification_token_expiry" pg:"verification_token_expiry"`
	TokenVersion            int       `json:"token_version" pg:"token_version"`
	CreatedAt               time.Time `json:"created_at" pg:"created_at"`
	UpdatedAt               time.Time `json:"updated_at" pg:"updated_at"`
	IsPremium               bool      `json:"is_premium" pg:"is_premium"`
	MaxUrls                 uint64    `json:"max_urls" pg:"max_urls"`
}

// AccessTokenClaims is a struct that contains the claims of the access token

/*
Function to generate JWT Token with user id, email and name
with given JWT_SECRET and JWT_EXPIRY
*/
func (user *User) GenerateJWTToken() (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":            user.Id,
		"email":         user.Email,
		"name":          user.Name,
		"token_version": user.TokenVersion,
		"isa":           time.Now().Unix(),
		"exp":           time.Now().Add(time.Hour * 24 * time.Duration(config.ACCESS_TOKEN_EXPIRY)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.ACCESS_TOKEN_SECRET))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

/*
Function to Update RefreshToken JWT Token with user id
*/
func (user *User) GenerateRefreshToken() (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.Id,
		"isa": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24 * time.Duration(config.REFRESH_TOKEN_EXPIRY)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.REFRESH_TOKEN_SECRET))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
