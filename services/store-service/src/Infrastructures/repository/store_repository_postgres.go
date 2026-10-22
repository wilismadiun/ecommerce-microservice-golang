package repository

import (
	"errors"
	"store-service/src/Domains/entities"

	"gorm.io/gorm"
)

type StoreRepositoryPostgres struct {
	DB *gorm.DB
}

func (h *StoreRepositoryPostgres) FindStoreByName(name string) (entities.Store, error) {
	var store entities.Store

	err := h.DB.Where("name = ?", name).First(&store).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Store{}, entities.ErrStoreNotFound
		}

		return entities.Store{}, err
	}

	return store, nil
}

func (h *StoreRepositoryPostgres) AddStore(store *entities.Store) error {
	err := h.DB.Create(store).Error
	if err != nil {
		return err
	}

	return nil
}

func (h *StoreRepositoryPostgres) DeleteStoreById(id string) error {
	store := entities.Store{}

	err := h.DB.Where("id = ?", id).Delete(&store).Error
	if err != nil {
		return err
	}

	return nil
}
