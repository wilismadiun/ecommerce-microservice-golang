package repository

import "order-service/src/Domains/entities"

type OrderRepository interface {
	CreateOrder(order entities.Order) error
}
