package repository

import (
	"order-service/src/Domains/entities"

	"gorm.io/gorm"
)

type OrderRepositoryPostgres struct {
	DB *gorm.DB
}

func (h *OrderRepositoryPostgres) CreateOrder(order entities.Order) error {
	err := h.DB.Create(&order).Error
	if err != nil {
		return entities.ErrDatabase
	}

	// items := order.Items
	// err = h.DB.Create(&items).Error
	// if err != nil {
	// 	h.DB.Delete(&order)

	// 	return entities.ErrDatabase
	// }

	return nil
}
