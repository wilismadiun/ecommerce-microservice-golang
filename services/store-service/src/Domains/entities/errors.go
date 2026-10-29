package entities

import "errors"

var (
	ErrUserNotFound       = errors.New("User Not Found")
	ErrStoreNotFound      = errors.New("Store Not Found")
	ErrStoreAlreadyExists = errors.New("Store Already Exist")
	ErrUserRoleNotFound   = errors.New("Access denied or data not found")
	ErrDatabase           = errors.New("Database operation failed")
)
