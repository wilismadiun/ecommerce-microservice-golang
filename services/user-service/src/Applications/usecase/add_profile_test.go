package usecase

import (
	"errors"
	"testing"
	"user-service/src/Domains/entities"
	"user-service/src/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_CreateProfile(t *testing.T) {
	t.Run("should return error when phone empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProfileRepo := mocks.NewMockProfileRepository(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)

		addProfile := &AddProfile{
			Repo:      mockProfileRepo,
			Generator: mockGenerator,
		}

		// Phone dikosongkan agar gagal pada tahap VerifyProfile
		inputProfile := entities.Profile{
			Name:   "John Doe",
			Phone:  "",
			UserId: "user-123",
		}

		res, err := addProfile.Execute(inputProfile)

		assert.Error(t, err)
		assert.Equal(t, "phone is required", err.Error())
		assert.Equal(t, entities.Profile{}, res)
	})

	t.Run("should return error when user id already exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProfileRepo := mocks.NewMockProfileRepository(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)

		mockGenerator.EXPECT().Generator().Return("profile-id-123").Times(1)

		mockProfileRepo.EXPECT().FindProfileByUserId("user-123").Return(entities.Profile{}, nil).Times(1)

		addProfile := &AddProfile{
			Repo:      mockProfileRepo,
			Generator: mockGenerator,
		}

		inputProfile := entities.Profile{
			Name:   "John Doe",
			Phone:  "081234567890",
			UserId: "user-123",
		}

		res, err := addProfile.Execute(inputProfile)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrProfileAlreadyExists, err)
		assert.Equal(t, entities.Profile{}, res)
	})

	t.Run("should return error when repository failed to create profile", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProfileRepo := mocks.NewMockProfileRepository(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)

		mockGenerator.EXPECT().Generator().Return("profile-id-123").Times(1)

		mockProfileRepo.EXPECT().FindProfileByUserId("user-123").Return(entities.Profile{}, entities.ErrUserNotFound).Times(1)

		repoErr := errors.New("database connection error")
		// Simulasi kegagalan saat insert ke database
		mockProfileRepo.EXPECT().CreateProfile(gomock.Any()).Return(repoErr).Times(1)

		addProfile := &AddProfile{
			Repo:      mockProfileRepo,
			Generator: mockGenerator,
		}

		inputProfile := entities.Profile{
			Name:   "John Doe",
			Phone:  "081234567890",
			UserId: "user-123",
		}

		res, err := addProfile.Execute(inputProfile)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		assert.Equal(t, entities.Profile{}, res)
	})

	t.Run("create profile success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProfileRepo := mocks.NewMockProfileRepository(ctrl)
		mockGenerator := mocks.NewMockGeneratorId(ctrl)

		mockGenerator.EXPECT().Generator().Return("profile-id-123").Times(1)

		mockProfileRepo.EXPECT().FindProfileByUserId("user-123").Return(entities.Profile{}, entities.ErrUserNotFound).Times(1)

		mockProfileRepo.EXPECT().CreateProfile(gomock.Any()).Return(nil).Times(1)

		addProfile := &AddProfile{
			Repo:      mockProfileRepo,
			Generator: mockGenerator,
		}

		inputProfile := entities.Profile{
			Name:   "John Doe",
			Phone:  "081234567890",
			UserId: "user-123",
		}

		expectedProfile := entities.Profile{
			ID:     "profile-id-123",
			Name:   "John Doe",
			Phone:  "081234567890",
			UserId: "user-123",
		}

		res, err := addProfile.Execute(inputProfile)

		assert.NoError(t, err)
		assert.Equal(t, expectedProfile, res)
	})
}
