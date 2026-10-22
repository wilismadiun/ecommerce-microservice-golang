package usecase

import (
	"user-service/src/Domains/entities"
	"user-service/src/Domains/repository"
)

type GetProfileByUser struct {
	Repo repository.ProfileRepository
}

func (h *GetProfileByUser) Execute(userId string) (entities.Profile, error) {
	profile, err := h.Repo.FindProfileByUserId(userId)
	if err != nil {
		return entities.Profile{}, err
	}

	return profile, nil
}
