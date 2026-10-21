package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var hasher = &HashPasswordBcrypt{}

func TestHashPasswordBcrypt_HashingPassword_Success(t *testing.T) {

	password := "password123"

	hash, err := hasher.Hash(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	err = bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	assert.NoError(t, err)
}

func Test_CompareHashPassword(t *testing.T) {
	password := "password-123"

	t.Run("should be error when password end password has do not match", func(t *testing.T) {
		err := hasher.CompareHashPassword(password, "newpass-123")

		assert.Error(t, err)
	})

	t.Run("compare success", func(t *testing.T) {
		hashPassword, err := hasher.Hash(password)

		assert.NoError(t, err)
		assert.NotEqual(t, password, hashPassword)
	})
}
