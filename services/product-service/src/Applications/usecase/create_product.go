package usecase

import (
	"product-service/src/Applications/generator"
	"product-service/src/Applications/storeclient"
	"product-service/src/Domains/entities"
	"product-service/src/Domains/repository"
)

type CreateProduct struct {
	ProductRepository repository.ProductRepository
	Generator         generator.GeneratorId
	StoreClient       storeclient.StoreClient
}

func (h *CreateProduct) Execute(token string, product entities.Product) (entities.Product, error) {
	roleUser, err := h.StoreClient.GetUserRole(token, product.StoreID)
	if err != nil {
		return entities.Product{}, err
	}

	if roleUser != "OWNER" && roleUser != "ADMIN" {
		return entities.Product{}, entities.ErrForbidden
	}

	err = entities.VerifyProdyct(product)
	if err != nil {
		return entities.Product{}, err
	}

	product.ID = h.Generator.Generator()

	err = h.ProductRepository.CreateProduct(&product)
	if err != nil {
		return entities.Product{}, err
	}

	return product, nil
}
