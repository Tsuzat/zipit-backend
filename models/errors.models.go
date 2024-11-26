package models

import "errors"

var (
	UserAlreadyExists  = errors.New("User already exists")
	UserNotFound       = errors.New("User not found")
	UserNotVerified    = errors.New("User not verified")
	EmailNotFoundError = errors.New("Email not found")
	PasswordIncorrect  = errors.New("Password incorrect")
)
