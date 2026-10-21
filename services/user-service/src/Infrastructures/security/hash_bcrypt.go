package security

import "golang.org/x/crypto/bcrypt"

type HashPasswordBcrypt struct{}

func (h *HashPasswordBcrypt) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hashedPassword := string(hash)

	return hashedPassword, nil
}

func (h *HashPasswordBcrypt) CompareHashPassword(password, hashPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
