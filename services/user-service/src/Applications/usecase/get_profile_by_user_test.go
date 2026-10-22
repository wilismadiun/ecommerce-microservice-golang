package usecase

import (
	"testing"
	"user-service/src/Domains/entities"
	"user-service/src/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_GetProfileByUser(t *testing.T) {
	t.Run("error when no profile found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProfileRepo := mocks.NewMockProfileRepository(ctrl)

		userId := "user-123"
		repoErr := entities.ErrProfileNotFound // atau error lain yang dikembalikan repository

		// Mock repository mengembalikan error saat profil tidak ditemukan
		mockProfileRepo.EXPECT().
			FindProfileByUserId(userId).
			Return(entities.Profile{}, repoErr).
			Times(1)

		useCase := &GetProfileByUser{
			Repo: mockProfileRepo,
		}

		profile, err := useCase.Execute(userId)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		assert.Equal(t, entities.Profile{}, profile)
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProfileRepo := mocks.NewMockProfileRepository(ctrl)

		userId := "user-123"
		expectedProfile := entities.Profile{
			ID:      "profile-001",
			UserId:  userId,
			Name:    "John Doe",
			Phone:   "081234567890",
			Address: "Jl. Merdeka No. 1",
		}

		// Mock repository berhasil menemukan dan mengembalikan profil
		mockProfileRepo.EXPECT().
			FindProfileByUserId(userId).
			Return(expectedProfile, nil).
			Times(1)

		useCase := &GetProfileByUser{
			Repo: mockProfileRepo,
		}

		profile, err := useCase.Execute(userId)

		assert.NoError(t, err)
		assert.Equal(t, expectedProfile, profile)
	})
}
