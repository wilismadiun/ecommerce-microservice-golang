package usecase

import (
	"errors"
	"testing"
	"user-service/src/Domains/entities"
	"user-service/src/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_register(t *testing.T) {
	t.Run("should return error when verify user fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockUserRepository(ctrl)
		mockHash := mocks.NewMockHasher(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)

		// Generator dipanggil di awal Execute sebelum VerifyUser dijalankan
		mockGenerator.EXPECT().Generator().Return("id-123").Times(1)

		u := &Register{
			Repo:      mockRepo,
			Hash:      mockHash,
			Generator: mockGenerator,
		}

		// Skenario: Email kosong
		inputUser := entities.User{
			Email:    "",
			Password: "validpassword123",
		}

		res, err := u.Execute(inputUser)

		// Validasi hasil
		assert.Error(t, err)
		assert.Equal(t, "email is required", err.Error())
		assert.Equal(t, entities.RegisteredUser{}, res)
	})

	t.Run("should return error when email already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockUserRepository(ctrl)
		mockHash := mocks.NewMockHasher(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)

		mockGenerator.EXPECT().Generator().Return("id-123").Times(1)
		// Simulasi email sudah terdaftar (err == nil)
		mockRepo.EXPECT().FindUserByEmail("test@example.com").Return(entities.User{}, nil).Times(1)

		u := &Register{
			Repo:      mockRepo,
			Hash:      mockHash,
			Generator: mockGenerator,
		}

		inputUser := entities.User{
			Email:    "test@example.com",
			Password: "validpassword123",
		}

		res, err := u.Execute(inputUser)

		assert.Error(t, err)
		assert.Equal(t, "Email is already in use", err.Error())
		assert.Equal(t, entities.RegisteredUser{}, res)
	})

	t.Run("should return error when find user returns unexpected error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockUserRepository(ctrl)
		mockHash := mocks.NewMockHasher(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)

		mockGenerator.EXPECT().Generator().Return("id-123").Times(1)

		unexpectedErr := errors.New("database connection failed")
		mockRepo.EXPECT().FindUserByEmail("test@example.com").Return(entities.User{}, unexpectedErr).Times(1)

		u := &Register{
			Repo:      mockRepo,
			Hash:      mockHash,
			Generator: mockGenerator,
		}

		inputUser := entities.User{
			Email:    "test@example.com",
			Password: "validpassword123",
		}

		res, err := u.Execute(inputUser)

		assert.Error(t, err)
		assert.Equal(t, unexpectedErr, err)
		assert.Equal(t, entities.RegisteredUser{}, res)
	})

	t.Run("should return error when hash password fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockUserRepository(ctrl)
		mockHash := mocks.NewMockHasher(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)

		mockGenerator.EXPECT().Generator().Return("id-123").Times(1)
		mockRepo.EXPECT().FindUserByEmail("test@example.com").Return(entities.User{}, entities.ErrUserNotFound).Times(1)

		hashErr := errors.New("failed to hash password")
		mockHash.EXPECT().Hash("validpassword123").Return("", hashErr).Times(1)

		u := &Register{
			Repo:      mockRepo,
			Hash:      mockHash,
			Generator: mockGenerator,
		}

		inputUser := entities.User{
			Email:    "test@example.com",
			Password: "validpassword123",
		}

		res, err := u.Execute(inputUser)

		assert.Error(t, err)
		assert.Equal(t, hashErr, err)
		assert.Equal(t, entities.RegisteredUser{}, res)
	})

	t.Run("should return error when repository failed to create user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockUserRepository(ctrl)
		mockHash := mocks.NewMockHasher(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)

		mockGenerator.EXPECT().Generator().Return("id-123").Times(1)
		mockRepo.EXPECT().FindUserByEmail("test@example.com").Return(entities.User{}, entities.ErrUserNotFound).Times(1)
		mockHash.EXPECT().Hash("validpassword123").Return("hashed_password_123", nil).Times(1)

		repoErr := errors.New("failed to insert user into database")
		mockRepo.EXPECT().UserRegister(gomock.Any()).Return(repoErr).Times(1)

		u := &Register{
			Repo:      mockRepo,
			Hash:      mockHash,
			Generator: mockGenerator,
		}

		inputUser := entities.User{
			Email:    "test@example.com",
			Password: "validpassword123",
		}

		res, err := u.Execute(inputUser)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		assert.Equal(t, entities.RegisteredUser{}, res)
	})

	t.Run("create user success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockUserRepository(ctrl)
		mockHash := mocks.NewMockHasher(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)

		mockGenerator.EXPECT().Generator().Return("id-123").Times(1)
		mockRepo.EXPECT().FindUserByEmail("test@example.com").Return(entities.User{}, entities.ErrUserNotFound).Times(1)
		mockHash.EXPECT().Hash("validpassword123").Return("hashed_password_123", nil).Times(1)
		mockRepo.EXPECT().UserRegister(gomock.Any()).Return(nil).Times(1)

		u := &Register{
			Repo:      mockRepo,
			Hash:      mockHash,
			Generator: mockGenerator,
		}

		inputUser := entities.User{
			Email:    "test@example.com",
			Password: "validpassword123",
		}

		res, err := u.Execute(inputUser)

		assert.NoError(t, err)
		assert.Equal(t, entities.RegisteredUser{
			ID:    "id-123",
			Email: "test@example.com",
		}, res)
	})
}
