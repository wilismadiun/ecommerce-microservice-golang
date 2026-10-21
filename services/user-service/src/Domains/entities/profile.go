package entities

import (
	"errors"
	"regexp"
	"strings"
)

type Profile struct {
	ID      string `json:"id"`
	UserId  string
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"` // Bersifat opsional (tidak wajib)
}

func VerifyProfile(profile Profile) error {
	if strings.TrimSpace(profile.Name) == "" {
		return errors.New("name is required")
	}

	if strings.TrimSpace(profile.Phone) == "" {
		return errors.New("phone is required")
	}

	phoneRegex := regexp.MustCompile(`^\d{10,15}$`)
	if !phoneRegex.MatchString(profile.Phone) {
		return errors.New("phone must be between 10 and 15 digits")
	}

	return nil
}
