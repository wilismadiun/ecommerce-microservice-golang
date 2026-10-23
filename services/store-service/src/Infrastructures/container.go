package infrastructures

import (
	"store-service/src/Applications/usecase"
	"store-service/src/Infrastructures/generator"
	"store-service/src/Infrastructures/repository"
	"store-service/src/Infrastructures/userclient"
	"store-service/src/Interfaces/http"

	"gorm.io/gorm"
)

func Container(db *gorm.DB) *http.Handler {
	storeRepoImplement := repository.StoreRepositoryPostgres{DB: db}
	userStoreRepoImplement := repository.UserStoreRepositoryPostgres{DB: db}
	generatorImplement := generator.GeneratorUUID{}
	UserCLientRequestImplement := userclient.UserCLientHttp{}

	addStoreHandler := usecase.AddStore{
		StoreRepository:     &storeRepoImplement,
		UserStoreRepository: &userStoreRepoImplement,
		Generator:           &generatorImplement,
		UserCLient:          &UserCLientRequestImplement,
	}

	return &http.Handler{
		AddStoreHandler: &addStoreHandler,
	}
}
