package repository

import (
	"errors"
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

func (h *UserStoreRepositoryPostgres) GetUserRole(userId, storeId string) (string, error) {
	var storeUser entities.UserStore

	err := h.DB.Preload("Role").Where("store_id = ? AND user_id = ?", storeId, userId).First(&storeUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", entities.ErrUserRoleNotFound
		}
		return "", entities.ErrDatabase
	}

	return storeUser.Role.Name, nil
}
