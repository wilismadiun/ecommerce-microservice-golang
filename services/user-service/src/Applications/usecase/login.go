package usecase

import (
	"user-service/src/Applications/security"
	"user-service/src/Domains/entities"
	"user-service/src/Domains/repository"
)

type Login struct {
	UserRepo repository.UserRepository
	Token    security.AuthToken
	Hasher   security.Hasher
}

func (h *Login) Execute(login entities.Login) (string, error) {
	err := entities.VerifyLogin(login)
	if err != nil {
		return "", err
	}

	existingUser, err := h.UserRepo.FindUserByEmail(login.Email)
	if err != nil {
		return "", err
	}

	err = h.Hasher.CompareHashPassword(login.Password, existingUser.Password)
	if err != nil {
		return "", err
	}

	accessToken, err := h.Token.GenerateToken(existingUser.ID)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
