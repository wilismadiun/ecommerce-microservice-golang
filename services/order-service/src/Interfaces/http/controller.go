package http

import (
	"errors"
	"net/http"
	"order-service/src/Applications/usecase"
	"order-service/src/Domains/entities"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	CreateOrderHandler *usecase.CreateOrder
}

func (h *Handler) CreateOrder(c *gin.Context) {
	userId := c.Param("id")

	var createOrder entities.Order
	err := c.ShouldBindBodyWithJSON(&createOrder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	createOrder.UserID = userId

	finishedOrder, err := h.CreateOrderHandler.Execute(createOrder)
	if err != nil {
		statusCode := http.StatusBadRequest

		if errors.Is(err, entities.ErrDatabase) {
			statusCode = http.StatusInternalServerError
		}

		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"meesage": "Order success",
		"data":    finishedOrder,
	})
}
