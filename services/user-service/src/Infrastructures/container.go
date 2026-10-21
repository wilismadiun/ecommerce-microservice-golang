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

	registerHandler := usecase.Register{
		Repo:      &userRepoImplement,
		Hash:      &hasherImplement,
		Generator: &generatorImplement,
	}

	return &http.Handler{
		RegisterHandler: &registerHandler,
	}
}
