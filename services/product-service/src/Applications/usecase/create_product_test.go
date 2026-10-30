package usecase

import (
	"errors"
	"product-service/src/Domains/entities"
	"product-service/src/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_CreateProduct(t *testing.T) {
	token := "Bearer secret-token-xyz" // Tambahkan variabel token untuk simulasi

	// Produk dasar yang valid untuk digunakan dalam pengujian
	validProduct := entities.Product{
		ID:          "prod-001",
		StoreID:     "store-999",
		Name:        "Mechanical Keyboard",
		Description: "Wireless RGB Mechanical Keyboard dengan switch Gateron",
		Attributes: map[string]interface{}{
			"brand":        "Keychron",
			"switch":       "Brown",
			"connectivity": "Bluetooth/USB-C",
			"weight_kg":    0.9,
		},
		Price: 1250000.0,
		Stock: 25,
	}

	t.Run("should return error when request user role fail or not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepo := mocks.NewMockProductRepository(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)
		mockStoreClient := mocks.NewMockStoreClient(ctrl)

		// Perbaikan: Menggunakan token, bukan userId
		mockStoreClient.EXPECT().
			GetUserRole(token, validProduct.StoreID).
			Return("", errors.New("service unavailable"))

		uc := &CreateProduct{
			ProductRepository: mockProductRepo,
			Generator:         mockGenerator,
			StoreClient:       mockStoreClient,
		}

		// Perbaikan: Menambahkan parameter token saat pemanggilan Execute
		result, err := uc.Execute(token, validProduct)

		assert.Error(t, err)
		assert.Equal(t, entities.Product{}, result)
	})

	t.Run("should return error when user role not admin or owner", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepo := mocks.NewMockProductRepository(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)
		mockStoreClient := mocks.NewMockStoreClient(ctrl)

		mockStoreClient.EXPECT().
			GetUserRole(token, validProduct.StoreID).
			Return("CASHIER", nil)

		uc := &CreateProduct{
			ProductRepository: mockProductRepo,
			Generator:         mockGenerator,
			StoreClient:       mockStoreClient,
		}

		result, err := uc.Execute(token, validProduct)

		assert.Error(t, err)
		assert.ErrorIs(t, err, entities.ErrForbidden)
		assert.Equal(t, entities.Product{}, result)
	})

	t.Run("should return error when name product empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepo := mocks.NewMockProductRepository(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)
		mockStoreClient := mocks.NewMockStoreClient(ctrl)

		mockStoreClient.EXPECT().
			GetUserRole(token, validProduct.StoreID).
			Return("OWNER", nil)

		invalidProduct := validProduct
		invalidProduct.Name = ""

		uc := &CreateProduct{
			ProductRepository: mockProductRepo,
			Generator:         mockGenerator,
			StoreClient:       mockStoreClient,
		}

		result, err := uc.Execute(token, invalidProduct)

		assert.Error(t, err)
		assert.Equal(t, entities.Product{}, result)
	})

	t.Run("should return error when repository failed to add product", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepo := mocks.NewMockProductRepository(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)
		mockStoreClient := mocks.NewMockStoreClient(ctrl)

		mockStoreClient.EXPECT().
			GetUserRole(token, validProduct.StoreID).
			Return("ADMIN", nil)

		generatedID := "prod-abc-123"
		mockGenerator.EXPECT().
			Generator().
			Return(generatedID)

		expectedProduct := validProduct
		expectedProduct.ID = generatedID
		mockProductRepo.EXPECT().
			CreateProduct(&expectedProduct).
			Return(errors.New("database connection error"))

		uc := &CreateProduct{
			ProductRepository: mockProductRepo,
			Generator:         mockGenerator,
			StoreClient:       mockStoreClient,
		}

		result, err := uc.Execute(token, validProduct)

		assert.Error(t, err)
		assert.Equal(t, entities.Product{}, result)
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepo := mocks.NewMockProductRepository(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)
		mockStoreClient := mocks.NewMockStoreClient(ctrl)

		mockStoreClient.EXPECT().
			GetUserRole(token, validProduct.StoreID).
			Return("OWNER", nil)

		generatedID := "prod-xyz-789"
		mockGenerator.EXPECT().
			Generator().
			Return(generatedID)

		expectedProduct := validProduct
		expectedProduct.ID = generatedID
		mockProductRepo.EXPECT().
			CreateProduct(&expectedProduct).
			Return(nil)

		uc := &CreateProduct{
			ProductRepository: mockProductRepo,
			Generator:         mockGenerator,
			StoreClient:       mockStoreClient,
		}

		result, err := uc.Execute(token, validProduct)

		assert.NoError(t, err)
		assert.Equal(t, generatedID, result.ID)
		assert.Equal(t, validProduct.Name, result.Name)
	})
}
