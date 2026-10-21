package usecase

import (
	"user-service/src/Applications/generator"
	"user-service/src/Domains/entities"
	"user-service/src/Domains/repository"
)

type AddProfile struct {
	Repo      repository.ProfileRepository
	Generator generator.GeneratorId
}

func (h *AddProfile) Execute(profile entities.Profile) (entities.Profile, error) {
	err := entities.VerifyProfile(profile)
	if err != nil {
		return entities.Profile{}, err
	}

	profile.ID = h.Generator.Generator()

	_, err = h.Repo.FindProfileByUserId(profile.UserId)
	if err == nil {
		return entities.Profile{}, entities.ErrProfileAlreadyExists
	}

	err = h.Repo.CreateProfile(&profile)
	if err != nil {
		return entities.Profile{}, err
	}

	return profile, nil
}
