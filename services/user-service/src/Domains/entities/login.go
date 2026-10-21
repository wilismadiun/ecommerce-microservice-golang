package entities

import (
	"errors"
	"regexp"
)

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func VerifyLogin(login Login) error {
	if login.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(
		`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`,
	)

	if !emailRegex.MatchString(login.Email) {
		return errors.New("invalid email format")
	}

	if login.Password == "" {
		return errors.New("password is required")
	}

	if len(login.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return nil
}
