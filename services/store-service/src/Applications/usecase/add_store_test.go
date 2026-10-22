package usecase

import (
	"errors"
	"store-service/src/Domains/entities"
	"store-service/src/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_AddStore(t *testing.T) {
	userId := "user-123"
	validStore := entities.Store{
		Name:        "Toko Maju",
		Description: "Toko kelontong",
		Phone:       "081234567890",
		Email:       "maju@example.com",
		Address:     "Jl. Sudirman",
	}

	t.Run("should return error when user request fail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStoreRepo := mocks.NewMockStoreRepository(ctrl)
		mockUserStoreRepo := mocks.NewMockUserStoreRepository(ctrl)
		mockUserClient := mocks.NewMockUserClient(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)

		useCase := &AddStore{
			StoreRepository:     mockStoreRepo,
			UserStoreRepository: mockUserStoreRepo,
			Generator:           mockGenerator,
			UserCLient:          mockUserClient,
		}

		mockUserClient.EXPECT().
			IsUserExists(userId).
			Return(false, errors.New("user service unavailable"))

		result, err := useCase.Execute(userId, validStore)

		assert.Error(t, err)
		assert.Equal(t, entities.Store{}, result)
		assert.Contains(t, err.Error(), "user service unavailable")
	})

	t.Run("should return error when email empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStoreRepo := mocks.NewMockStoreRepository(ctrl)
		mockUserStoreRepo := mocks.NewMockUserStoreRepository(ctrl)
		mockUserClient := mocks.NewMockUserClient(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)

		useCase := &AddStore{
			StoreRepository:     mockStoreRepo,
			UserStoreRepository: mockUserStoreRepo,
			Generator:           mockGenerator,
			UserCLient:          mockUserClient,
		}

		invalidStore := validStore
		invalidStore.Email = "" // Mengosongkan email agar gagal di VerifyStore

		mockUserClient.EXPECT().
			IsUserExists(userId).
			Return(true, nil)

		result, err := useCase.Execute(userId, invalidStore)

		assert.Error(t, err)
		assert.Equal(t, entities.Store{}, result)
		assert.Equal(t, "email is required", err.Error())
	})

	t.Run("should return error when name store already exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStoreRepo := mocks.NewMockStoreRepository(ctrl)
		mockUserStoreRepo := mocks.NewMockUserStoreRepository(ctrl)
		mockUserClient := mocks.NewMockUserClient(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)

		useCase := &AddStore{
			StoreRepository:     mockStoreRepo,
			UserStoreRepository: mockUserStoreRepo,
			Generator:           mockGenerator,
			UserCLient:          mockUserClient,
		}

		mockUserClient.EXPECT().
			IsUserExists(userId).
			Return(true, nil)

		mockStoreRepo.EXPECT().
			FindStoreByName(validStore.Name).
			Return(entities.Store{}, errors.New("Store found"))

		result, err := useCase.Execute(userId, validStore)

		assert.Error(t, err)
		assert.Equal(t, entities.Store{}, result)
	})

	t.Run("should return error when user has 3 stores", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStoreRepo := mocks.NewMockStoreRepository(ctrl)
		mockUserStoreRepo := mocks.NewMockUserStoreRepository(ctrl)
		mockUserClient := mocks.NewMockUserClient(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)

		useCase := &AddStore{
			StoreRepository:     mockStoreRepo,
			UserStoreRepository: mockUserStoreRepo,
			Generator:           mockGenerator,
			UserCLient:          mockUserClient,
		}

		mockUserClient.EXPECT().
			IsUserExists(userId).
			Return(true, nil)

		mockStoreRepo.EXPECT().
			FindStoreByName(validStore.Name).
			Return(entities.Store{}, entities.ErrStoreNotFound)

		mockUserStoreRepo.EXPECT().
			FindUserStoreByUser(userId).
			Return([]entities.UserStore{
				{Id: "1", StoreId: "s-1", UserId: userId, RoleId: "1"},
				{Id: "2", StoreId: "s-2", UserId: userId, RoleId: "1"},
				{Id: "3", StoreId: "s-3", UserId: userId, RoleId: "1"},
			}, nil)

		result, err := useCase.Execute(userId, validStore)

		assert.Error(t, err)
		assert.Equal(t, entities.Store{}, result)
		assert.Equal(t, "Max store limit reached", err.Error())
	})

	t.Run("should return error when repository failed to add store", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStoreRepo := mocks.NewMockStoreRepository(ctrl)
		mockUserStoreRepo := mocks.NewMockUserStoreRepository(ctrl)
		mockUserClient := mocks.NewMockUserClient(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)

		useCase := &AddStore{
			StoreRepository:     mockStoreRepo,
			UserStoreRepository: mockUserStoreRepo,
			Generator:           mockGenerator,
			UserCLient:          mockUserClient,
		}

		mockUserClient.EXPECT().
			IsUserExists(userId).
			Return(true, nil)

		mockStoreRepo.EXPECT().
			FindStoreByName(validStore.Name).
			Return(entities.Store{}, entities.ErrStoreNotFound)

		mockUserStoreRepo.EXPECT().
			FindUserStoreByUser(userId).
			Return([]entities.UserStore{}, nil)

		mockGenerator.EXPECT().
			Generator().
			Return("store-uuid-123")

		mockStoreRepo.EXPECT().
			AddStore(gomock.Any()).
			Return(errors.New("database connection error"))

		result, err := useCase.Execute(userId, validStore)

		assert.Error(t, err)
		assert.Equal(t, entities.Store{}, result)
		assert.Contains(t, err.Error(), "database connection error")
	})

	// error
	t.Run("should return error when repository failed to add owner", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStoreRepo := mocks.NewMockStoreRepository(ctrl)
		mockUserStoreRepo := mocks.NewMockUserStoreRepository(ctrl)
		mockUserClient := mocks.NewMockUserClient(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)

		useCase := &AddStore{
			StoreRepository:     mockStoreRepo,
			UserStoreRepository: mockUserStoreRepo,
			Generator:           mockGenerator,
			UserCLient:          mockUserClient,
		}

		mockUserClient.EXPECT().
			IsUserExists(userId).
			Return(true, nil)

		mockStoreRepo.EXPECT().
			FindStoreByName(validStore.Name).
			Return(entities.Store{}, entities.ErrStoreNotFound)

		mockUserStoreRepo.EXPECT().
			FindUserStoreByUser(userId).
			Return([]entities.UserStore{}, nil)

		mockGenerator.EXPECT().
			Generator().
			Return("store-uuid-123")

		mockStoreRepo.EXPECT().
			AddStore(gomock.Any()).
			Return(nil)

		mockGenerator.EXPECT().
			Generator().
			Return("user-store-uuid-123")

		mockUserStoreRepo.EXPECT().
			AddUserStore(gomock.Any()).
			Return(errors.New("failed to add owner"))

		mockStoreRepo.EXPECT().
			DeleteStoreById("store-uuid-123").
			Return(nil)

		result, err := useCase.Execute(userId, validStore)

		assert.Error(t, err)
		assert.Equal(t, entities.Store{}, result)
		assert.Equal(t, "failed to add owner", err.Error())
	})

	// error
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStoreRepo := mocks.NewMockStoreRepository(ctrl)
		mockUserStoreRepo := mocks.NewMockUserStoreRepository(ctrl)
		mockUserClient := mocks.NewMockUserClient(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)

		useCase := &AddStore{
			StoreRepository:     mockStoreRepo,
			UserStoreRepository: mockUserStoreRepo,
			Generator:           mockGenerator,
			UserCLient:          mockUserClient,
		}

		mockUserClient.EXPECT().
			IsUserExists(userId).
			Return(true, nil)

		mockStoreRepo.EXPECT().
			FindStoreByName(validStore.Name).
			Return(entities.Store{}, entities.ErrStoreNotFound)

		mockUserStoreRepo.EXPECT().
			FindUserStoreByUser(userId).
			Return([]entities.UserStore{}, nil)

		mockGenerator.EXPECT().
			Generator().
			Return("store-uuid-123")

		mockStoreRepo.EXPECT().
			AddStore(gomock.Any()).
			Return(nil)

		mockGenerator.EXPECT().
			Generator().
			Return("user-store-uuid-123")

		mockUserStoreRepo.EXPECT().
			AddUserStore(gomock.Any()).
			Return(nil)

		expectedStore := validStore
		expectedStore.Id = "store-uuid-123"

		result, err := useCase.Execute(userId, validStore)

		assert.NoError(t, err)
		assert.Equal(t, expectedStore, result)
	})
}
