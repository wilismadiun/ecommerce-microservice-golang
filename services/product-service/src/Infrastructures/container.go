package infrastructures

import (
	"product-service/src/Applications/usecase"
	"product-service/src/Infrastructures/generator"
	"product-service/src/Infrastructures/repository"
	"product-service/src/Infrastructures/storeclient"
	"product-service/src/Interfaces/http"

	"gorm.io/gorm"
)

func Container(db *gorm.DB) *http.Handler {
	productRepositoryImplement := repository.ProductRepositoryPostgres{DB: db}
	generatorImplement := generator.GeneratorUUID{}
	storeClientImplement := storeclient.StoreClientHttp{}

	createProductHandler := usecase.CreateProduct{
		ProductRepository: &productRepositoryImplement,
		Generator:         &generatorImplement,
		StoreClient:       &storeClientImplement,
	}

	decreaseStockhandler := usecase.DecreaseStock{
		ProductRepository: &productRepositoryImplement,
	}

	updateStockHandler := usecase.UpdateProduct{
		ProductRepository: &productRepositoryImplement,
		StoreClient:       &storeClientImplement,
	}

	return &http.Handler{
		CreateProductHanderl: &createProductHandler,
		DecreaseStockhandler: &decreaseStockhandler,
		UpdateStockHandler:   &updateStockHandler,
	}
}
