package usecase

import (
	"errors"
	"order-service/src/Domains/entities"
	"order-service/src/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_CreateOrder(t *testing.T) {
	t.Run("should return error when quantity empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockOrderRepository := mocks.NewMockOrderRepository(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)
		mockProductClient := mocks.NewMockProductClient(ctrl)

		uc := &CreateOrder{
			OrderRepository: mockOrderRepository,
			Generator:       mockGenerator,
			ProductClient:   mockProductClient,
		}

		inputOrder := entities.Order{
			StoreID:         "store-test-999",
			ShippingAddress: "Jl. Merdeka No. 45",
			Items: []entities.OrderItem{
				{
					ProductID:   "prod-001",
					ProductName: "Mechanical Keyboard",
					Quantity:    0, // Quantity invalid
					Price:       1250000.0,
				},
			},
		}

		result, err := uc.Execute(inputOrder)

		assert.Error(t, err)
		assert.Equal(t, entities.Order{}, result)
		assert.Equal(t, "quantity must be greater than zero", err.Error())
	})

	t.Run("should return error when stock is less than quantity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockOrderRepository := mocks.NewMockOrderRepository(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)
		mockProductClient := mocks.NewMockProductClient(ctrl)

		uc := &CreateOrder{
			OrderRepository: mockOrderRepository,
			Generator:       mockGenerator,
			ProductClient:   mockProductClient,
		}

		inputOrder := entities.Order{
			StoreID:         "store-test-999",
			ShippingAddress: "Jl. Merdeka No. 45",
			Items: []entities.OrderItem{
				{
					ProductID:   "prod-001",
					ProductName: "Mechanical Keyboard",
					Quantity:    5,
					Price:       1250000.0,
				},
			},
		}

		mockGenerator.EXPECT().
			Generator().
			Return("order-id-123")

		expectedErr := errors.New("insufficient stock")
		mockProductClient.EXPECT().
			DecreaseStock("prod-001", 5).
			Return(expectedErr)

		result, err := uc.Execute(inputOrder)

		assert.Error(t, err)
		assert.Equal(t, entities.Order{}, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should return error when repository failed to create order", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockOrderRepository := mocks.NewMockOrderRepository(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)
		mockProductClient := mocks.NewMockProductClient(ctrl)

		uc := &CreateOrder{
			OrderRepository: mockOrderRepository,
			Generator:       mockGenerator,
			ProductClient:   mockProductClient,
		}

		inputOrder := entities.Order{
			StoreID:         "store-test-999",
			ShippingAddress: "Jl. Merdeka No. 45",
			Items: []entities.OrderItem{
				{
					ProductID:   "prod-001",
					ProductName: "Mechanical Keyboard",
					Quantity:    2,
					Price:       1250000.0,
				},
			},
		}

		gomock.InOrder(
			mockGenerator.EXPECT().Generator().Return("order-id-123"),
			mockGenerator.EXPECT().Generator().Return("item-id-001"),
		)

		mockProductClient.EXPECT().
			DecreaseStock("prod-001", 2).
			Return(nil)

		mockOrderRepository.EXPECT().
			CreateOrder(gomock.Any()).
			Return(entities.ErrDatabase)

		result, err := uc.Execute(inputOrder)

		assert.Error(t, err)
		assert.Equal(t, entities.Order{}, result)
		assert.Equal(t, entities.ErrDatabase, err)
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockOrderRepository := mocks.NewMockOrderRepository(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)
		mockProductClient := mocks.NewMockProductClient(ctrl)

		uc := &CreateOrder{
			OrderRepository: mockOrderRepository,
			Generator:       mockGenerator,
			ProductClient:   mockProductClient,
		}

		inputOrder := entities.Order{
			StoreID:         "store-test-999",
			ShippingAddress: "Jl. Merdeka No. 45",
			Items: []entities.OrderItem{
				{
					ProductID:   "prod-001",
					ProductName: "Mechanical Keyboard",
					Quantity:    2,
					Price:       1250000.0,
				},
			},
		}

		gomock.InOrder(
			mockGenerator.EXPECT().Generator().Return("order-id-123"),
			mockGenerator.EXPECT().Generator().Return("item-id-001"),
		)

		mockProductClient.EXPECT().
			DecreaseStock("prod-001", 2).
			Return(nil)

		mockOrderRepository.EXPECT().
			CreateOrder(gomock.Any()).
			Return(nil)

		result, err := uc.Execute(inputOrder)

		assert.NoError(t, err)
		assert.Equal(t, "order-id-123", result.ID)
		assert.Equal(t, "PENDING", result.Status)
		assert.Equal(t, 2500000.0, result.TotalAmount)
		assert.Len(t, result.Items, 1)
		assert.Equal(t, "item-id-001", result.Items[0].ID)
		assert.Equal(t, "order-id-123", result.Items[0].OrderID)
	})
}
