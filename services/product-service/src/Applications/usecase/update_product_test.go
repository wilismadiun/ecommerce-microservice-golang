package usecase

import (
	"errors"
	"product-service/src/Domains/entities"
	"product-service/src/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_UpdateStock(t *testing.T) {
	token := "valid-token-123"
	productID := "prod-001"
	storeID := "store-999"

	t.Run("should return error when product id empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepository := mocks.NewMockProductRepository(ctrl)
		mockStoreClient := mocks.NewMockStoreClient(ctrl)

		uc := &UpdateProduct{
			ProductRepository: mockProductRepository,
			StoreClient:       mockStoreClient,
		}

		// Input dengan ProductID kosong (akan langsung gagal di VerifyUpdateStock)
		inputData := entities.UpdateStock{
			ProductID: "",
			Quantity:  10,
		}

		err := uc.Execute(token, inputData)

		assert.Error(t, err)
		assert.Equal(t, "product id must be available", err.Error())
	})

	t.Run("should return error when product not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepository := mocks.NewMockProductRepository(ctrl)
		mockStoreClient := mocks.NewMockStoreClient(ctrl)

		inputData := entities.UpdateStock{
			ProductID: productID,
			Quantity:  10,
		}

		// Mock GetProductById mengembalikan error (misal record not found)
		expectedErr := errors.New("product not found")
		mockProductRepository.EXPECT().
			GetProductById(productID).
			Return(entities.Product{}, expectedErr)

		uc := &UpdateProduct{
			ProductRepository: mockProductRepository,
			StoreClient:       mockStoreClient,
		}

		err := uc.Execute(token, inputData)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should return error when request user role fail or not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepository := mocks.NewMockProductRepository(ctrl)
		mockStoreClient := mocks.NewMockStoreClient(ctrl)

		inputData := entities.UpdateStock{
			ProductID: productID,
			Quantity:  10,
		}

		// Mock GetProductById sukses dan mengembalikan data produk dengan StoreID tertentu
		mockProductRepository.EXPECT().
			GetProductById(productID).
			Return(entities.Product{ID: productID, StoreID: storeID}, nil)

		// Mock GetUserRole gagal / mengembalikan error
		expectedErr := errors.New("failed to connect to store service")
		mockStoreClient.EXPECT().
			GetUserRole(token, storeID).
			Return("", expectedErr)

		uc := &UpdateProduct{
			ProductRepository: mockProductRepository,
			StoreClient:       mockStoreClient,
		}

		err := uc.Execute(token, inputData)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should return error when user role not admin or owner", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepository := mocks.NewMockProductRepository(ctrl)
		mockStoreClient := mocks.NewMockStoreClient(ctrl)

		inputData := entities.UpdateStock{
			ProductID: productID,
			Quantity:  10,
		}

		mockProductRepository.EXPECT().
			GetProductById(productID).
			Return(entities.Product{ID: productID, StoreID: storeID}, nil)

		// Mock GetUserRole mengembalikan role selain OWNER atau ADMIN (misal: "CASHIER" atau "")
		mockStoreClient.EXPECT().
			GetUserRole(token, storeID).
			Return("CASHIER", nil)

		uc := &UpdateProduct{
			ProductRepository: mockProductRepository,
			StoreClient:       mockStoreClient,
		}

		err := uc.Execute(token, inputData)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrForbidden, err)
	})

	t.Run("should return error when repository failed to update product", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepository := mocks.NewMockProductRepository(ctrl)
		mockStoreClient := mocks.NewMockStoreClient(ctrl)

		inputData := entities.UpdateStock{
			ProductID: productID,
			Quantity:  10,
		}

		mockProductRepository.EXPECT().
			GetProductById(productID).
			Return(entities.Product{ID: productID, StoreID: storeID}, nil)

		// Mock GetUserRole mengembalikan role yang diizinkan (OWNER)
		mockStoreClient.EXPECT().
			GetUserRole(token, storeID).
			Return("OWNER", nil)

		// Mock UpdateStock gagal / mengembalikan error database
		expectedErr := errors.New("database update error")
		mockProductRepository.EXPECT().
			UpdateStock(inputData.ProductID, inputData.Quantity).
			Return(expectedErr)

		uc := &UpdateProduct{
			ProductRepository: mockProductRepository,
			StoreClient:       mockStoreClient,
		}

		err := uc.Execute(token, inputData)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepository := mocks.NewMockProductRepository(ctrl)
		mockStoreClient := mocks.NewMockStoreClient(ctrl)

		inputData := entities.UpdateStock{
			ProductID: productID,
			Quantity:  10,
		}

		mockProductRepository.EXPECT().
			GetProductById(productID).
			Return(entities.Product{ID: productID, StoreID: storeID}, nil)

		// Mock GetUserRole mengembalikan role ADMIN yang valid
		mockStoreClient.EXPECT().
			GetUserRole(token, storeID).
			Return("ADMIN", nil)

		// Mock UpdateStock berhasil (return nil)
		mockProductRepository.EXPECT().
			UpdateStock(inputData.ProductID, inputData.Quantity).
			Return(nil)

		uc := &UpdateProduct{
			ProductRepository: mockProductRepository,
			StoreClient:       mockStoreClient,
		}

		err := uc.Execute(token, inputData)

		assert.NoError(t, err)
	})
}
