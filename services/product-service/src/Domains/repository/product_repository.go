package repository

import "product-service/src/Domains/entities"

type ProductRepository interface {
	CreateProduct(product *entities.Product) error
}
