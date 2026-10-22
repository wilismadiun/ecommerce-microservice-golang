package entities

import (
	"testing"
)

func TestVerifyStore(t *testing.T) {
	tests := []struct {
		name    string
		store   Store
		wantErr bool
	}{
		{
			name: "Valid Store (Tanpa ID)",
			store: Store{
				Name:        "Toko Berkah",
				Description: "Menjual berbagai kebutuhan pokok",
				Phone:       "081234567890",
				Email:       "berkah@example.com",
				Address:     "Jl. Merdeka No. 10",
			},
			wantErr: false,
		},
		{
			name: "Valid Store (Dengan ID)",
			store: Store{
				Id:          "store-uuid-123",
				Name:        "Toko Berkah",
				Description: "Menjual berbagai kebutuhan pokok",
				Phone:       "081234567890",
				Email:       "berkah@example.com",
				Address:     "Jl. Merdeka No. 10",
			},
			wantErr: false,
		},
		{
			name: "Error: Missing Name",
			store: Store{
				Description: "Menjual berbagai kebutuhan pokok",
				Phone:       "081234567890",
				Email:       "berkah@example.com",
				Address:     "Jl. Merdeka No. 10",
			},
			wantErr: true,
		},
		{
			name: "Error: Missing Description",
			store: Store{
				Name:    "Toko Berkah",
				Phone:   "081234567890",
				Email:   "berkah@example.com",
				Address: "Jl. Merdeka No. 10",
			},
			wantErr: true,
		},
		{
			name: "Error: Missing Phone",
			store: Store{
				Name:        "Toko Berkah",
				Description: "Menjual berbagai kebutuhan pokok",
				Email:       "berkah@example.com",
				Address:     "Jl. Merdeka No. 10",
			},
			wantErr: true,
		},
		{
			name: "Error: Phone Terlalu Pendek (< 10)",
			store: Store{
				Name:        "Toko Berkah",
				Description: "Menjual berbagai kebutuhan pokok",
				Phone:       "123456789", // 9 karakter
				Email:       "berkah@example.com",
				Address:     "Jl. Merdeka No. 10",
			},
			wantErr: true,
		},
		{
			name: "Error: Phone Terlalu Panjang (> 15)",
			store: Store{
				Name:        "Toko Berkah",
				Description: "Menjual berbagai kebutuhan pokok",
				Phone:       "12345678901234567", // 17 karakter
				Email:       "berkah@example.com",
				Address:     "Jl. Merdeka No. 10",
			},
			wantErr: true,
		},
		{
			name: "Error: Missing Email",
			store: Store{
				Name:        "Toko Berkah",
				Description: "Menjual berbagai kebutuhan pokok",
				Phone:       "081234567890",
				Address:     "Jl. Merdeka No. 10",
			},
			wantErr: true,
		},
		{
			name: "Error: Format Email Salah",
			store: Store{
				Name:        "Toko Berkah",
				Description: "Menjual berbagai kebutuhan pokok",
				Phone:       "081234567890",
				Email:       "berkah-at-example.com",
				Address:     "Jl. Merdeka No. 10",
			},
			wantErr: true,
		},
		{
			name: "Error: Missing Address",
			store: Store{
				Name:        "Toko Berkah",
				Description: "Menjual berbagai kebutuhan pokok",
				Phone:       "081234567890",
				Email:       "berkah@example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyStore(tt.store)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyStore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
