package repository

import (
	"store-service/src/Domains/entities"

	"gorm.io/gorm"
)

type UserStoreRepositoryPostgres struct {
	DB *gorm.DB
}

func (h *UserStoreRepositoryPostgres) FindUserStoreByUser(userId string) ([]entities.UserStore, error) {
	var usersStore []entities.UserStore

	err := h.DB.Where("user_id = ?", userId).Find(&usersStore).Error
	if err != nil {
		return nil, err
	}

	return usersStore, nil
}

func (h *UserStoreRepositoryPostgres) AddUserStore(userStore entities.UserStore) error {
	err := h.DB.Create(&userStore).Error
	if err != nil {
		return err
	}

	return nil
}
