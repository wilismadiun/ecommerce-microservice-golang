package repository

import (
	"user-service/src/Domains/entities"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (h *UserRepository) FindUserByEmail(email string) (entities.User, error) {
	var user entities.User

	err := h.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entities.User{}, entities.ErrNotFound
		}

		return entities.User{}, err
	}

	return user, nil
}

func (h *UserRepository) UserRegister(user *entities.User) error {
	err := h.DB.Create(user).Error
	if err != nil {
		return err
	}

	return nil
}
