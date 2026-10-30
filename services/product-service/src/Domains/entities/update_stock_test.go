package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_VerifyUpdateStock(t *testing.T) {
	testCases := []struct {
		name          string
		input         UpdateStock
		expectedError string
	}{
		{
			name: "product id empty",
			input: UpdateStock{
				ProductID: "",
				Quantity:  10,
			},
			expectedError: "product id must be available",
		},
		{
			name: "quantity < 0",
			input: UpdateStock{
				ProductID: "prod-001",
				Quantity:  -5,
			},
			expectedError: "stock quantity cannot be negative",
		},
		{
			name: "success",
			input: UpdateStock{
				ProductID: "prod-001",
				Quantity:  0,
			},
			expectedError: "",
		},
	}

	// Melakukan perulangan untuk menjalankan setiap subtest
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := VerifyUpdateStock(tc.input)

			if tc.expectedError != "" {
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tc.expectedError, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
