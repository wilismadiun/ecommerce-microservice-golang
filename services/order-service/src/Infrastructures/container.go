package infrastructures

import (
	"order-service/src/Applications/usecase"
	"order-service/src/Infrastructures/generator"
	"order-service/src/Infrastructures/productclient"
	"order-service/src/Infrastructures/repository"
	"order-service/src/Interfaces/http"

	"gorm.io/gorm"
)

func Container(db *gorm.DB) *http.Handler {
	orderRepositoryImplement := repository.OrderRepositoryPostgres{DB: db}
	generatorImplement := generator.GeneratorUUID{}
	productClientImplement := productclient.ProductClientHttp{}

	createOrderHandler := usecase.CreateOrder{
		OrderRepository: &orderRepositoryImplement,
		Generator:       &generatorImplement,
		ProductClient:   &productClientImplement,
	}

	return &http.Handler{
		CreateOrderHandler: &createOrderHandler,
	}
}
