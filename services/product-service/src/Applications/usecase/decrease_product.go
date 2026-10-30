package usecase

import (
	"product-service/src/Domains/entities"
	"product-service/src/Domains/repository"
)

type DecreaseStock struct {
	ProductRepository repository.ProductRepository
}

func (h *DecreaseStock) Execute(decreaseStock entities.DecreaseStock) error {
	err := entities.VerifyDecreaseProduct(decreaseStock)
	if err != nil {
		return err
	}

	err = h.ProductRepository.DecreaseStock(decreaseStock.ProductID, decreaseStock.Quantity)
	if err != nil {
		return err
	}

	return nil
}
