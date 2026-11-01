package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyCreateOrder(t *testing.T) {
	validItem := OrderItem{
		ProductID:   "prod-001",
		ProductName: "Mechanical Keyboard",
		Quantity:    2,
		Price:       1250000.0,
	}

	testCases := []struct {
		name          string
		input         Order
		expectedError string
	}{
		{
			name: "store id is empty",
			input: Order{
				StoreID:         "",
				ShippingAddress: "Jl. Sudirman No. 123",
				Items:           []OrderItem{validItem},
			},
			expectedError: "store_id is required",
		},
		{
			name: "shipping address is empty",
			input: Order{
				StoreID:         "store-test-999",
				ShippingAddress: "",
				Items:           []OrderItem{validItem},
			},
			expectedError: "shipping_address is required",
		},
		{
			name: "items list is empty",
			input: Order{
				StoreID:         "store-test-999",
				ShippingAddress: "Jl. Sudirman No. 123",
				Items:           []OrderItem{},
			},
			expectedError: "order must contain at least one item",
		},
		{
			name: "item product id is empty",
			input: Order{
				StoreID:         "store-test-999",
				ShippingAddress: "Jl. Sudirman No. 123",
				Items: []OrderItem{
					{
						ProductID:   "",
						ProductName: "Mechanical Keyboard",
						Quantity:    2,
						Price:       1250000.0,
					},
				},
			},
			expectedError: "product_id is required for all items",
		},
		{
			name: "item product name is empty",
			input: Order{
				StoreID:         "store-test-999",
				ShippingAddress: "Jl. Sudirman No. 123",
				Items: []OrderItem{
					{
						ProductID:   "prod-001",
						ProductName: "",
						Quantity:    2,
						Price:       1250000.0,
					},
				},
			},
			expectedError: "Product name is required for all items",
		},
		{
			name: "item quantity is less than or equal to zero",
			input: Order{
				StoreID:         "store-test-999",
				ShippingAddress: "Jl. Sudirman No. 123",
				Items: []OrderItem{
					{
						ProductID:   "prod-001",
						ProductName: "Mechanical Keyboard",
						Quantity:    0, // atau negatif seperti -1
						Price:       1250000.0,
					},
				},
			},
			expectedError: "quantity must be greater than zero",
		},
		{
			name: "item price is negative",
			input: Order{
				StoreID:         "store-test-999",
				ShippingAddress: "Jl. Sudirman No. 123",
				Items: []OrderItem{
					{
						ProductID:   "prod-001",
						ProductName: "Mechanical Keyboard",
						Quantity:    2,
						Price:       -500.0, // harga negatif
					},
				},
			},
			expectedError: "price cannot be negative",
		},
		{
			name: "success",
			input: Order{
				StoreID:         "store-test-999",
				ShippingAddress: "Jl. Sudirman No. 123",
				Items:           []OrderItem{validItem},
			},
			expectedError: "",
		},
	}

	// Menjalankan setiap skenario uji menggunakan perulangan subtest
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := VerifyCreateOrder(tc.input)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
