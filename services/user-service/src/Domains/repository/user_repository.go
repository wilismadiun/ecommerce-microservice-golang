package repository

import "user-service/src/Domains/entities"

type UserRepository interface {
	UserRegister(*entities.User) error
	FindUserByEmail(email string) (entities.User, error)
}
