package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyDecreaseProduct(t *testing.T) {
	// Definisikan struct untuk skenario test case
	testCases := []struct {
		name          string
		input         DecreaseStock
		expectedError string
	}{
		{
			name: "product id tidak ada",
			input: DecreaseStock{
				ProductID: "",
				Quantity:  5,
			},
			expectedError: "The product id must be available",
		},
		{
			name: "quantity <= 0 (nol)",
			input: DecreaseStock{
				ProductID: "prod-001",
				Quantity:  0,
			},
			expectedError: "quantity must be greater than 0",
		},
		{
			name: "quantity <= 0 (negatif)",
			input: DecreaseStock{
				ProductID: "prod-001",
				Quantity:  -2,
			},
			expectedError: "quantity must be greater than 0",
		},
		{
			name: "sukses",
			input: DecreaseStock{
				ProductID: "prod-001",
				Quantity:  10,
			},
			expectedError: "",
		},
	}

	// Lakukan perulangan untuk menjalankan setiap sub test
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := VerifyDecreaseProduct(tc.input)

			if tc.expectedError != "" {
				// Jika diharapkan ada error
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				// Jika diharapkan sukses (tidak ada error)
				assert.NoError(t, err)
			}
		})
	}
}
