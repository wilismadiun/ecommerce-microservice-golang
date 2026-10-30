package entities

import "errors"

var (
	ErrForbidden        = errors.New("You do not have permission to access this resource")
	ErrDatabase         = errors.New("Database operation failed")
	ErrUserRoleNotFound = errors.New("user role not found")
)
