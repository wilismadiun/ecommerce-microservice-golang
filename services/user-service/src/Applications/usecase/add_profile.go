package usecase

import "user-service/src/Domains/repository"

type AddProfile struct {
	Repo repository.ProfileRepository
}
