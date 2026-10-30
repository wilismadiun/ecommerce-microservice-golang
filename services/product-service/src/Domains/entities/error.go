package entities

import "errors"

var (
	ErrForbidden         = errors.New("You do not have permission to access this resource")
	ErrDatabase          = errors.New("Database operation failed")
	ErrUserRoleNotFound  = errors.New("user role not found")
	ErrInsufficientStock = errors.New("product stock is not enough")
	ErrProductNotFound   = errors.New("Product not found")
)
