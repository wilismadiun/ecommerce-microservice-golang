package repository

import (
	"errors"
	"user-service/src/Domains/entities"

	"gorm.io/gorm"
)

type ProfileRepository struct {
	DB *gorm.DB
}

func (h *ProfileRepository) CreateProfile(profile *entities.Profile) error {
	err := h.DB.Create(profile).Error
	if err != nil {
		return err
	}

	return nil
}

func (h *ProfileRepository) FindProfileByUserId(userId string) (entities.Profile, error) {
	var profile entities.Profile

	err := h.DB.Where("user_id = ?", userId).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Profile{}, entities.ErrProfileNotFound
		}

		return entities.Profile{}, err
	}

	return profile, nil
}
