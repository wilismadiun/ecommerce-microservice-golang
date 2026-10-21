package entities

import (
	"errors"
	"regexp"
	"strings"
)

type Profile struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"` // Bersifat opsional (tidak wajib)
}

func VerifyProfile(profile Profile) error {
	// 1. Validasi Nama Wajib Ada
	if strings.TrimSpace(profile.Name) == "" {
		return errors.New("name is required")
	}

	// 2. Validasi Phone Wajib Ada
	if strings.TrimSpace(profile.Phone) == "" {
		return errors.New("phone is required")
	}

	// 3. Validasi Phone harus berupa angka dengan panjang 10 - 15 karakter
	phoneRegex := regexp.MustCompile(`^\d{10,15}$`)
	if !phoneRegex.MatchString(profile.Phone) {
		return errors.New("phone must be between 10 and 15 digits")
	}

	// 4. Address bersifat opsional, jadi tidak ada blok validasi error jika kosong.
	return nil
}
