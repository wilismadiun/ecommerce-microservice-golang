package usecase

import (
	"errors"
	"user-service/src/Applications/generator"
	"user-service/src/Applications/security"
	"user-service/src/Domains/entities"
	"user-service/src/Domains/repository"
)

type Register struct {
	Repo      repository.UserRepository
	Hash      security.Hasher
	Generator generator.GeneratorId
}

func (h *Register) Execute(user entities.User) (entities.RegisteredUser, error) {
	user.ID = h.Generator.Generator()

	err := entities.VerifyUser(user)
	if err != nil {
		return entities.RegisteredUser{}, err
	}

	_, err = h.Repo.FindUserByEmail(user.Email)
	if err == nil {
		return entities.RegisteredUser{}, errors.New("Email is already in use")
	}

	if !errors.Is(err, entities.ErrNotFound) {
		return entities.RegisteredUser{}, err
	}

	hashedPassword, err := h.Hash.Hash(user.Password)
	if err != nil {
		return entities.RegisteredUser{}, err
	}

	user.Password = hashedPassword

	err = h.Repo.UserRegister(&user)
	if err != nil {
		return entities.RegisteredUser{}, err
	}

	return entities.RegisteredUser{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}
