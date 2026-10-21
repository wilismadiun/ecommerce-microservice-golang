package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyLogin(t *testing.T) {
	t.Run("should return error when email is empty", func(t *testing.T) {
		login := Login{
			Email:    "",
			Password: "password123",
		}

		err := VerifyLogin(login)

		assert.Error(t, err)
		assert.Equal(t, "email is required", err.Error())
	})

	t.Run("should return error when email format is invalid", func(t *testing.T) {
		login := Login{
			Email:    "invalid-email",
			Password: "password123",
		}

		err := VerifyLogin(login)

		assert.Error(t, err)
		assert.Equal(t, "invalid email format", err.Error())
	})

	t.Run("should return error when password is empty", func(t *testing.T) {
		login := Login{
			Email:    "user@example.com",
			Password: "",
		}

		err := VerifyLogin(login)

		assert.Error(t, err)
		assert.Equal(t, "password is required", err.Error())
	})

	t.Run("should return error when password is less than 8 characters", func(t *testing.T) {
		login := Login{
			Email:    "user@example.com",
			Password: "1234567",
		}

		err := VerifyLogin(login)

		assert.Error(t, err)
		assert.Equal(t, "password must be at least 8 characters", err.Error())
	})
}
