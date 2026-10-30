package repository

import "product-service/src/Domains/entities"

type ProductRepository interface {
	CreateProduct(product *entities.Product) error
	DecreaseStock(productID string, quantity int) error
	GetProductById(id string) (entities.Product, error)
	UpdateStock(productId string, quantity int) error
}
