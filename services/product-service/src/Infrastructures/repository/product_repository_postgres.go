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

func (h *ProductRepositoryPostgres) DecreaseStock(productID string, quantity int) error {
	result := h.DB.Table("products").Where("id = ? AND stock >= ?", productID, quantity).Update("stock", gorm.Expr("stock - ?", quantity))

	err := result.Error
	if err != nil {
		return entities.ErrDatabase
	}

	if result.RowsAffected == 0 {
		return entities.ErrInsufficientStock
	}

	return nil
}
