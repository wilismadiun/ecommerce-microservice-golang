package repository

import "user-service/src/Domains/entities"

type ProfileRepository interface {
	CreateProfile(*entities.Profile) error
	FindProfileByUserId(userId string) (entities.Profile, error)
}
