package usecase

import (
	"errors"
	"store-service/src/mocks"
	"testing"

	"go.uber.org/mock/gomock"
)

func Test_GetUserRole(t *testing.T) {
	t.Run("should return error when user not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserStoreRepository := mocks.NewMockUserStoreRepository(ctrl)

		userId := "user-123"
		storeId := "store-123"
		expectedErr := errors.New("user not found")

		// Menyiapkan ekspektasi mock repository
		mockUserStoreRepository.EXPECT().
			GetUserRole(userId, storeId).
			Return("", expectedErr)

		// Inisialisasi usecase dengan mock repository
		useCase := &GetUserRole{
			UserStoreRepository: mockUserStoreRepository,
		}

		// Eksekusi fungsi
		role, err := useCase.Execute(userId, storeId)

		// Assertions
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != expectedErr.Error() {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
		if role != "" {
			t.Fatalf("expected empty role, got %v", role)
		}
	})

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserStoreRepository := mocks.NewMockUserStoreRepository(ctrl)

		userId := "user-123"
		storeId := "store-123"
		expectedRole := "ADMIN"

		// Menyiapkan ekspektasi mock repository
		mockUserStoreRepository.EXPECT().
			GetUserRole(userId, storeId).
			Return(expectedRole, nil)

		// Inisialisasi usecase dengan mock repository
		useCase := &GetUserRole{
			UserStoreRepository: mockUserStoreRepository,
		}

		// Eksekusi fungsi
		role, err := useCase.Execute(userId, storeId)

		// Assertions
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if role != expectedRole {
			t.Fatalf("expected role %v, got %v", expectedRole, role)
		}
	})
}
