package entities

import "errors"

type Order struct {
	ID              string      `gorm:"primaryKey" json:"id"`
	UserID          string      `json:"user_id"`
	StoreID         string      `json:"store_id"`
	TotalAmount     float64     `json:"total_amount"`
	Status          string      `json:"status"`
	ShippingAddress string      `json:"shipping_address"`
	Items           []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
}

type OrderItem struct {
	ID          string `gorm:"primaryKey" json:"id"`
	OrderID     string
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

func VerifyCreateOrder(req Order) error {
	if req.StoreID == "" {
		return errors.New("store_id is required")
	}
	if req.ShippingAddress == "" {
		return errors.New("shipping_address is required")
	}
	if len(req.Items) == 0 {
		return errors.New("order must contain at least one item")
	}

	for _, item := range req.Items {
		if item.ProductID == "" {
			return errors.New("product_id is required for all items")
		}
		if item.ProductName == "" {
			return errors.New("Product name is required for all items")
		}
		if item.Quantity <= 0 {
			return errors.New("quantity must be greater than zero")
		}
		if item.Price < 0 {
			return errors.New("price cannot be negative")
		}
	}

	return nil
}
