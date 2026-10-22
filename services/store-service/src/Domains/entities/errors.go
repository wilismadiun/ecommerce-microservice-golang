package entities

import "errors"

var (
	ErrUserNotFound       = errors.New("User Not Found")
	ErrStoreNotFound      = errors.New("Store Not Found")
	ErrStoreAlreadyExists = errors.New("Store Already Exist")
)
