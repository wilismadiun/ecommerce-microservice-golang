package repository

import (
	"product-service/src/Domains/entities"

	"gorm.io/gorm"
)

type ProductRepositoryPostgres struct {
	DB *gorm.DB
}

func (h *ProductRepositoryPostgres) CreateProduct(product *entities.Product) error {
	err := h.DB.Create(product).Error
	if err != nil {
		return entities.ErrDatabase
	}

	return nil
}
