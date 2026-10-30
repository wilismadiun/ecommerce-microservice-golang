package usecase

import (
	"errors"
	"product-service/src/Domains/entities"
	"product-service/src/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_DecreaseStock(t *testing.T) {
	productId := "prod-001"
	quantity := 5

	decreaseStock := entities.DecreaseStock{
		ProductID: productId,
		Quantity:  quantity,
	}

	t.Run("should return error when product id empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockProductRepository := mocks.NewMockProductRepository(ctrl)

		uc := &DecreaseStock{
			ProductRepository: mockProductRepository,
		}

		input := entities.DecreaseStock{
			ProductID: "",
			Quantity:  quantity,
		}

		err := uc.Execute(input)

		assert.Error(t, err)
		assert.Equal(t, "The product id must be available", err.Error())
	})

	t.Run("should return error when quantity invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockProductRepository := mocks.NewMockProductRepository(ctrl)

		uc := &DecreaseStock{
			ProductRepository: mockProductRepository,
		}

		input := entities.DecreaseStock{
			ProductID: productId,
			Quantity:  0,
		}

		err := uc.Execute(input)

		assert.Error(t, err)
	})

	t.Run("should return error when database error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockProductRepository := mocks.NewMockProductRepository(ctrl)

		expectedErr := errors.New("database connection error")

		mockProductRepository.EXPECT().
			DecreaseStock(productId, quantity).
			Return(expectedErr)

		uc := &DecreaseStock{
			ProductRepository: mockProductRepository,
		}

		err := uc.Execute(decreaseStock)

		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("should return error when insufficient stock", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockProductRepository := mocks.NewMockProductRepository(ctrl)

		mockProductRepository.EXPECT().
			DecreaseStock(productId, quantity).
			Return(entities.ErrInsufficientStock)

		uc := &DecreaseStock{
			ProductRepository: mockProductRepository,
		}

		err := uc.Execute(decreaseStock)

		assert.ErrorIs(t, err, entities.ErrInsufficientStock)
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockProductRepository := mocks.NewMockProductRepository(ctrl)

		mockProductRepository.EXPECT().
			DecreaseStock(productId, quantity).
			Return(nil)

		uc := &DecreaseStock{
			ProductRepository: mockProductRepository,
		}

		err := uc.Execute(decreaseStock)

		assert.NoError(t, err)
	})
}
