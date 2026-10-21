package entities

import "errors"

var (
	ErrNotFound        = errors.New("User Not Found")
	ErrInvalidPassword = errors.New("invalid password")
)
