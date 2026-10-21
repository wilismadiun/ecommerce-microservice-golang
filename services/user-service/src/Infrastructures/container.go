package infrastructures

import (
	"user-service/src/Applications/usecase"
	"user-service/src/Infrastructures/generator"
	"user-service/src/Infrastructures/repository"
	"user-service/src/Infrastructures/security"
	"user-service/src/Interfaces/http"

	"gorm.io/gorm"
)

func Container(db *gorm.DB) *http.Handler {
	userRepoImplement := repository.UserRepository{DB: db}
	hasherImplement := security.HashPasswordBcrypt{}
	generatorImplement := generator.GeneratorUUID{}
	tokenImplement := security.AuthenticationTokenJWT{}

	registerHandler := usecase.Register{
		Repo:      &userRepoImplement,
		Hash:      &hasherImplement,
		Generator: &generatorImplement,
	}

	loginHandler := usecase.Login{
		UserRepo: &userRepoImplement,
		Token:    &tokenImplement,
		Hasher:   &hasherImplement,
	}

	return &http.Handler{
		RegisterHandler: &registerHandler,
		LoginHandler:    &loginHandler,
	}
}
