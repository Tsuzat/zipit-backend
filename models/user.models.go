package models

import (
	"time"

	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Id                      string    `json:"id" pg:"id,pk,type:uuid,default:gen_random_uuid(),notnull"`
	Name                    string    `json:"name" pg:"name,notnull"`
	Email                   string    `json:"email" pg:"email,unique,notnull"`
	Password                string    `json:"password" pg:"password,notnull"`
	ProfileImage            string    `json:"profile_image" pg:"profile_image"`
	IsVerified              bool      `json:"is_verified" pg:"is_verified,default:false"`
	RefreshToken            string    `json:"refresh_token" pg:"refresh_token"`
	RefreshTokenExpiry      time.Time `json:"refresh_token_expiry" pg:"refresh_token_expiry"`
	VerificationToken       string    `json:"verification_token" pg:"verification_token"`
	VerificationTokenExpiry time.Time `json:"verification_token_expiry" pg:"verification_token_expiry"`
	CreatedAt               time.Time `json:"created_at" pg:"created_at,default:now()"`
	UpdatedAt               time.Time `json:"updated_at" pg:"updated_at,default:now()"`
	IsPremium               bool      `json:"is_premium" pg:"is_premium,default:false"`
	MaxUrls                 uint64    `json:"max_urls" pg:"max_urls,default:10"`
}

/*
Function to generate JWT Token with user id, email and name
with given JWT_SECRET and JWT_EXPIRY
*/
func (user *User) GenerateJWTToken() (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"name":  user.Name,
		"isa":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 24 * time.Duration(config.ACCESS_TOKEN_EXPIRY)).Unix(),
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
