package repository

import "store-service/src/Domains/entities"

type UserStoreRepository interface {
	FindUserStoreByUser(userId string) ([]entities.UserStore, error)
	AddUserStore(userStore entities.UserStore) error
	GetUserRole(userId, storeId string) (string, error)
}
