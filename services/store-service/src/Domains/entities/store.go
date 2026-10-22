package entities

import (
	"errors"
	"regexp"
	"strings"
)

type Store struct {
	Id          string
	Name        string `json:"name"`
	Description string `json:"description"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Address     string `json:"address"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func VerifyStore(store Store) error {
	// 1. Validasi semua field wajib diisi kecuali Id
	if strings.TrimSpace(store.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(store.Description) == "" {
		return errors.New("description is required")
	}
	if strings.TrimSpace(store.Phone) == "" {
		return errors.New("phone is required")
	}
	if strings.TrimSpace(store.Email) == "" {
		return errors.New("email is required")
	}
	if strings.TrimSpace(store.Address) == "" {
		return errors.New("address is required")
	}

	// 2. Validasi panjang phone harus antara 10 - 15 karakter
	if len(store.Phone) < 10 || len(store.Phone) > 15 {
		return errors.New("phone must be between 10 and 15 characters")
	}

	// 3. Validasi format email
	if !emailRegex.MatchString(store.Email) {
		return errors.New("invalid email format")
	}

	return nil
}
