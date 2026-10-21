package entities

import "errors"

var (
	ErrUserNotFound         = errors.New("User Not Found")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrProfileAlreadyExists = errors.New("Profile already exists")
	ErrProfileNotFound      = errors.New("Profile Not Found")
)
