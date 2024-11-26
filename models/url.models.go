package models

import (
	"time"
)

type Url struct {
	Id        string    `json:"id" pg:"id,pk,type:uuid,default:gen_random_uuid(),notnull"`
	Url       string    `json:"url" pg:"url,notnull"`
	Alias     string    `json:"alias" pg:"alias,notnull,unique"`
	CreatedAt time.Time `json:"created_at" pg:"created_at,default:now()"`
	UpdatedAt time.Time `json:"updated_at" pg:"updated_at,default:now()"`
	ExpiresAt time.Time `json:"expires_at" pg:"expires_at,notnull"`
	OwnerId   string    `json:"owner_id" pg:"owner_id,type:uuid,notnull"`
	Owner     *User     `json:"owner" pg:"fk:owner_id,rel:has-one,on_delete:cascade"`
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
