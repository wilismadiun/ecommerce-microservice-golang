package entities_test

import (
	"product-service/src/Domains/entities"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyProdyct(t *testing.T) {
	// Base data produk yang valid
	baseProduct := entities.Product{
		ID:          "prod-001",
		StoreID:     "store-999",
		Name:        "Mechanical Keyboard",
		Description: "Wireless RGB Mechanical Keyboard dengan switch Gateron",
		Attributes: map[string]interface{}{
			"brand":        "Keychron",
			"switch":       "Brown",
			"connectivity": "Bluetooth/USB-C",
			"weight_kg":    0.9,
		},
		Price: 1250000.0,
		Stock: 25,
	}

	// Definisi test cases dalam bentuk tabel
	testCases := []struct {
		name        string
		modify      func(p *entities.Product)
		expectError bool
	}{
		{
			name: "1. Store ID tidak ada",
			modify: func(p *entities.Product) {
				p.StoreID = ""
			},
			expectError: true,
		},
		{
			name: "2. Name tidak ada",
			modify: func(p *entities.Product) {
				p.Name = ""
			},
			expectError: true,
		},
		{
			name: "3. Description tidak ada",
			modify: func(p *entities.Product) {
				p.Description = ""
			},
			expectError: true,
		},
		{
			name: "4. Price tidak ada (0)",
			modify: func(p *entities.Product) {
				p.Price = 0
			},
			expectError: true,
		},
		{
			name: "5. Stock tidak ada (0)",
			modify: func(p *entities.Product) {
				p.Stock = 0
			},
			expectError: true,
		},
		{
			name: "6. Price negatif",
			modify: func(p *entities.Product) {
				p.Price = -150000.0
			},
			expectError: true,
		},
		{
			name: "7. Stock negatif",
			modify: func(p *entities.Product) {
				p.Stock = -5
			},
			expectError: true,
		},
		{
			name: "8. Sukses dengan atribut lengkap",
			modify: func(p *entities.Product) {
				// Tidak ada yang diubah, gunakan data valid murni
			},
			expectError: false,
		},
	}

	// Menjalankan test cases menggunakan for loop dan subtest
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Salin data baseProduct agar modifikasi tidak mencemari test case lain
			p := baseProduct
			tc.modify(&p)

			err := entities.VerifyProdyct(p)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, p.Attributes)
				assert.Equal(t, "Keychron", p.Attributes["brand"])
			}
		})
	}
}
