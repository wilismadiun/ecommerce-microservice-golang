package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyUser(t *testing.T) {
	// Definisikan test cases dalam bentuk struct slice
	testCases := []struct {
		name          string
		user          User
		expectedError string
	}{
		{
			name: "Email kosong",
			user: User{
				ID:       "ID-001",
				Email:    "",
				Password: "validpassword123",
			},
			expectedError: "email is required",
		},
		{
			name: "Password kosong",
			user: User{
				ID:       "ID-001",
				Email:    "test@example.com",
				Password: "",
			},
			expectedError: "password is required",
		},
		{
			name: "Email tidak sesuai format",
			user: User{
				ID:       "ID-001",
				Email:    "invalid-email-format",
				Password: "validpassword123",
			},
			expectedError: "invalid email format",
		},
		{
			name: "Password kurang dari 8 karakter",
			user: User{
				ID:       "ID-001",
				Email:    "test@example.com",
				Password: "123",
			},
			expectedError: "password must be at least 8 characters",
		},
	}

	// Loop untuk menjalankan setiap skenario test menggunakan t.Run
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := VerifyUser(tc.user)

			// Menggunakan testify/assert untuk validasi error
			assert.Error(t, err)
			assert.Equal(t, tc.expectedError, err.Error())
		})
	}
}
