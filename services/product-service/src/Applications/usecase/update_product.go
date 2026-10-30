package usecase

import (
	"product-service/src/Applications/storeclient"
	"product-service/src/Domains/entities"
	"product-service/src/Domains/repository"
)

type UpdateProduct struct {
	ProductRepository repository.ProductRepository
	StoreClient       storeclient.StoreClient
}

func (h *UpdateProduct) Execute(token string, data entities.UpdateStock) error {
	err := entities.VerifyUpdateStock(data)
	if err != nil {
		return err
	}

	existProduct, err := h.ProductRepository.GetProductById(data.ProductID)
	if err != nil {
		return err
	}

	roleUser, err := h.StoreClient.GetUserRole(token, existProduct.StoreID)
	if err != nil {
		return err
	}

	if roleUser != "OWNER" && roleUser != "ADMIN" {
		return entities.ErrForbidden
	}

	err = h.ProductRepository.UpdateStock(data.ProductID, data.Quantity)
	if err != nil {
		return err
	}

	return nil
}
