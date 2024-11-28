package models

import (
	"time"
)

type Url struct {
	Id        int       `json:"id" pg:"id"`
	Url       string    `json:"url" pg:"url"`
	Alias     string    `json:"alias" pg:"alias"`
	CreatedAt time.Time `json:"created_at" pg:"created_at"`
	UpdatedAt time.Time `json:"updated_at" pg:"updated_at"`
	ExpiresAt time.Time `json:"expires_at" pg:"expires_at"`
	Owner     int       `json:"owner" pg:"owner"`
}

type CreateUrlRequest struct {
	Url   string `json:"url"`
	Alias string `json:"alias"`
}

type UpdateUrlRequest struct {
	Url       string    `json:"url"`
	Alias     string    `json:"alias"`
	ExpiresAt time.Time `json:"expires_at"`
}
