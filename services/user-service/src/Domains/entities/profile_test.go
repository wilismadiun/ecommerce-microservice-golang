package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyProfile(t *testing.T) {
	testCases := []struct {
		name          string
		profile       Profile
		expectedError string
	}{
		{
			name: "Nama kosong",
			profile: Profile{
				ID:      "ID-001",
				UserId:  "user-123",
				Name:    "",
				Phone:   "081234567890",
				Address: "Jl. Merdeka No. 1",
			},
			expectedError: "name is required",
		},
		{
			name: "Phone kosong",
			profile: Profile{
				ID:      "ID-001",
				UserId:  "user-123",
				Name:    "John Doe",
				Phone:   "",
				Address: "Jl. Merdeka No. 1",
			},
			expectedError: "phone is required",
		},
		{
			name: "Phone kurang dari 10 digit",
			profile: Profile{
				ID:      "ID-001",
				UserId:  "user-123",
				Name:    "John Doe",
				Phone:   "0812345",
				Address: "Jl. Merdeka No. 1",
			},
			expectedError: "phone must be between 10 and 15 digits",
		},
		{
			name: "Phone lebih dari 15 digit",
			profile: Profile{
				ID:      "ID-001",
				UserId:  "user-123",
				Name:    "John Doe",
				Phone:   "081234567890123456",
				Address: "Jl. Merdeka No. 1",
			},
			expectedError: "phone must be between 10 and 15 digits",
		},
		{
			name: "Phone mengandung huruf/simbol",
			profile: Profile{
				ID:      "ID-001",
				UserId:  "user-123",
				Name:    "John Doe",
				Phone:   "0812345abcde",
				Address: "Jl. Merdeka No. 1",
			},
			expectedError: "phone must be between 10 and 15 digits",
		},
		{
			name: "Profile valid dengan Address diisi",
			profile: Profile{
				ID:      "ID-001",
				UserId:  "user-123",
				Name:    "John Doe",
				Phone:   "081234567890",
				Address: "Jl. Merdeka No. 1",
			},
			expectedError: "",
		},
		{
			name: "Profile valid dengan Address kosong (opsional)",
			profile: Profile{
				ID:      "ID-001",
				UserId:  "user-123",
				Name:    "John Doe",
				Phone:   "081234567890",
				Address: "", // Kosong tetap dianggap valid
			},
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		tc := tc // PENTING: Mencegah variabel loop closure pada goroutine/sub-test
		t.Run(tc.name, func(t *testing.T) {
			err := VerifyProfile(tc.profile)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
