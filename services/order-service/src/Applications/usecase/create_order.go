package usecase

import (
	"order-service/src/Applications/generator"
	"order-service/src/Applications/productclient"
	"order-service/src/Domains/entities"
	"order-service/src/Domains/repository"
)

type CreateOrder struct {
	OrderRepository repository.OrderRepository
	Generator       generator.GeneratorId
	ProductClient   productclient.ProductClient
}

func (h *CreateOrder) Execute(order entities.Order) (entities.Order, error) {
	err := entities.VerifyCreateOrder(order)
	if err != nil {
		return entities.Order{}, err
	}

	order.ID = h.Generator.Generator()

	var totalPrice float64 = 0.0

	for i, item := range order.Items {
		err := h.ProductClient.DecreaseStock(item.ProductID, item.Quantity)
		if err != nil {
			return entities.Order{}, err
		}

		price := item.Price * float64(item.Quantity)
		totalPrice += price

		order.Items[i].ID = h.Generator.Generator()
		order.Items[i].OrderID = order.ID
	}

	order.Status = "PENDING"
	order.TotalAmount = totalPrice

	err = h.OrderRepository.CreateOrder(order)
	if err != nil {
		return entities.Order{}, err
	}

	return order, nil
}
