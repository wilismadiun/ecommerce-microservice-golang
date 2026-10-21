package usecase

import (
	"testing"
	"user-service/src/Domains/entities"
	"user-service/src/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_Login(t *testing.T) {
	t.Run("should return error when email is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockToken := mocks.NewMockAuthToken(ctrl)
		mockHasher := mocks.NewMockHasher(ctrl)

		loginUseCase := &Login{
			UserRepo: mockUserRepo,
			Token:    mockToken,
			Hasher:   mockHasher,
		}

		input := entities.Login{
			Email:    "",
			Password: "validpassword123",
		}

		token, err := loginUseCase.Execute(input)

		assert.Error(t, err)
		assert.Equal(t, "", token)
	})

	t.Run("should return error when email is not registered", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockToken := mocks.NewMockAuthToken(ctrl)
		mockHasher := mocks.NewMockHasher(ctrl)

		input := entities.Login{
			Email:    "notfound@example.com",
			Password: "validpassword123",
		}

		// Mock repository mengembalikan error karena email tidak ditemukan
		mockUserRepo.EXPECT().
			FindUserByEmail(input.Email).
			Return(entities.User{}, entities.ErrUserNotFound).
			Times(1)

		loginUseCase := &Login{
			UserRepo: mockUserRepo,
			Token:    mockToken,
			Hasher:   mockHasher,
		}

		token, err := loginUseCase.Execute(input)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrUserNotFound, err)
		assert.Equal(t, "", token)
	})

	t.Run("should return error when password is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockToken := mocks.NewMockAuthToken(ctrl)
		mockHasher := mocks.NewMockHasher(ctrl)

		loginUseCase := &Login{
			UserRepo: mockUserRepo,
			Token:    mockToken,
			Hasher:   mockHasher,
		}

		input := entities.Login{
			Email:    "test@example.com",
			Password: "",
		}

		token, err := loginUseCase.Execute(input)

		assert.Error(t, err)
		assert.Equal(t, "", token)
	})

	t.Run("should return error when password does not match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockToken := mocks.NewMockAuthToken(ctrl)
		mockHasher := mocks.NewMockHasher(ctrl)

		input := entities.Login{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		existingUser := entities.User{
			ID:       "id-123",
			Email:    "test@example.com",
			Password: "hashed_correct_password",
		}

		// Mock repository mengembalikan user yang ditemukan
		mockUserRepo.EXPECT().
			FindUserByEmail(input.Email).
			Return(existingUser, nil).
			Times(1)

		// Mock hasher gagal mencocokkan password
		mockHasher.EXPECT().
			CompareHashPassword(input.Password, existingUser.Password).
			Return(entities.ErrInvalidPassword).
			Times(1)

		loginUseCase := &Login{
			UserRepo: mockUserRepo,
			Token:    mockToken,
			Hasher:   mockHasher,
		}

		token, err := loginUseCase.Execute(input)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrInvalidPassword, err)
		assert.Equal(t, "", token)
	})

	t.Run("should return access token when login is successful", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := mocks.NewMockUserRepository(ctrl)
		mockToken := mocks.NewMockAuthToken(ctrl)
		mockHasher := mocks.NewMockHasher(ctrl)

		input := entities.Login{
			Email:    "test@example.com",
			Password: "validpassword123",
		}

		existingUser := entities.User{
			ID:       "id-123",
			Email:    "test@example.com",
			Password: "hashed_correct_password",
		}

		// 1. Mock repository menemukan user
		mockUserRepo.EXPECT().
			FindUserByEmail(input.Email).
			Return(existingUser, nil).
			Times(1)

		// 2. Mock hasher berhasil mencocokkan password (return nil)
		mockHasher.EXPECT().
			CompareHashPassword(input.Password, existingUser.Password).
			Return(nil).
			Times(1)

		expectedToken := "mock-jwt-access-token-xyz"
		// 3. Mock token berhasil men-generate token berdasarkan ID user
		mockToken.EXPECT().
			GenerateToken(existingUser.ID).
			Return(expectedToken, nil).
			Times(1)

		loginUseCase := &Login{
			UserRepo: mockUserRepo,
			Token:    mockToken,
			Hasher:   mockHasher,
		}

		token, err := loginUseCase.Execute(input)

		assert.NoError(t, err)
		assert.Equal(t, expectedToken, token)
	})
}
