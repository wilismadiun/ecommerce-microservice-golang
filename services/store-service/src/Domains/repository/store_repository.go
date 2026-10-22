package repository

import "store-service/src/Domains/entities"

type StoreRepository interface {
	FindStoreByName(name string) (entities.Store, error)
	AddStore(store *entities.Store) error
	DeleteStoreById(id string) error
}
