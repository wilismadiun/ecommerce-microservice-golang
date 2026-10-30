package entities

import "errors"

type UpdateStock struct {
	ProductID string
	Quantity  int
}

func VerifyUpdateStock(data UpdateStock) error {
	if data.ProductID == "" {
		return errors.New("product id must be available")
	}

	if data.Quantity < 0 {
		return errors.New("stock quantity cannot be negative")
	}

	return nil
}
