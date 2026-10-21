package entities

import (
	"errors"
	"regexp"
	"strings"
)

type User struct {
	ID       string
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisteredUser struct {
	ID    string
	Email string
}

func VerifyUser(user User) error {
	if strings.TrimSpace(user.ID) == "" {
		return errors.New("id is required")
	}

	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is required")
	}

	if strings.TrimSpace(user.Password) == "" {
		return errors.New("password is required")
	}

	emailRegex := regexp.MustCompile(
		`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`,
	)

	if !emailRegex.MatchString(user.Email) {
		return errors.New("invalid email format")
	}

	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return nil
}
