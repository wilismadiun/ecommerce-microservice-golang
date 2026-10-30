package http

import (
	"errors"
	"net/http"
	"product-service/src/Applications/usecase"
	"product-service/src/Domains/entities"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	CreateProductHanderl *usecase.CreateProduct
}

func (h *Handler) CreateProduct(c *gin.Context) {
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized: user context not found",
		})
		return
	}

	userCtx, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to parse user context",
		})
		return
	}

	var product entities.Product
	err := c.ShouldBindBodyWithJSON(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	product, err = h.CreateProductHanderl.Execute(userCtx, product)
	if err != nil {
		statusCode := http.StatusBadRequest

		if errors.Is(err, entities.ErrForbidden) {
			statusCode = http.StatusForbidden
		}

		if errors.Is(err, entities.ErrDatabase) {
			statusCode = http.StatusInternalServerError
		}

		if errors.Is(err, entities.ErrUserRoleNotFound) {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product creation was successful",
		"data":    product,
	})
}
