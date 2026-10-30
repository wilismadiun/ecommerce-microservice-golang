package entities

import "errors"

type DecreaseStock struct {
	ProductID string
	Quantity  int
}

func VerifyDecreaseProduct(product DecreaseStock) error {
	if product.ProductID == "" {
		return errors.New("The product id must be available")
	}

	if product.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	return nil
}
