package repository

import (
	"log"
	"os"
	"testing"
	"user-service/config"
	"user-service/src/Domains/entities"
	"user-service/src/Infrastructures/database"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var userRepo *UserRepository
var profileRepo *ProfileRepository

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Println("Peringatan: Gagal memuat .env dari root, mencoba path alternatif...")
	}

	dsn := config.DSN()
	database.ConnectDB(dsn)

	userRepo = &UserRepository{
		DB: database.DB,
	}

	profileRepo = &ProfileRepository{
		DB: database.DB,
	}

	code := m.Run()

	os.Exit(code)
}

func Test_FindUserByEmail(t *testing.T) {
	t.Run("Email not found", func(t *testing.T) {
		email := "nonexistent@example.com"

		user, err := userRepo.FindUserByEmail(email)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrUserNotFound, err)
		assert.Equal(t, entities.User{}, user)
	})

	t.Run("Email found", func(t *testing.T) {
		email := "found@example.com"
		dummyUser := entities.User{
			ID:       "id-test-find-001",
			Email:    email,
			Password: "hashedpassword123",
		}

		// Masukkan data dummy ke database untuk pengujian
		err := userRepo.DB.Create(&dummyUser).Error
		assert.NoError(t, err)

		// Hapus data secara otomatis di akhir test (cleanup)
		t.Cleanup(func() {
			userRepo.DB.Where("email = ?", email).Delete(&entities.User{})
		})

		user, err := userRepo.FindUserByEmail(email)

		assert.NoError(t, err)
		assert.Equal(t, dummyUser.ID, user.ID)
		assert.Equal(t, dummyUser.Email, user.Email)
	})
}

func Test_UserRegister(t *testing.T) {
	t.Run("User register success", func(t *testing.T) {
		email := "register-success@gmail.com"
		newUser := &entities.User{
			ID:       "id-test-register-001",
			Email:    email,
			Password: "hashedpassword123",
		}

		// Hapus data secara otomatis di akhir test (cleanup)
		t.Cleanup(func() {
			userRepo.DB.Where("email = ?", email).Delete(&entities.User{})
		})

		err := userRepo.UserRegister(newUser)

		assert.NoError(t, err)

		// Validasi tambahan: Pastikan data benar-benar tersimpan di database
		var savedUser entities.User
		err = userRepo.DB.Where("email = ?", email).First(&savedUser).Error
		assert.NoError(t, err)
		assert.Equal(t, newUser.ID, savedUser.ID)
		assert.Equal(t, newUser.Email, savedUser.Email)
	})
}

func Test_FindProfileByUserId(t *testing.T) {
	t.Run("Profile not found", func(t *testing.T) {
		userId := "nonexistent-user-id"

		profile, err := profileRepo.FindProfileByUserId(userId)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrProfileNotFound, err)
		assert.Equal(t, entities.Profile{}, profile)
	})

	t.Run("Profile found", func(t *testing.T) {
		userId := "id-test-user-profile-001"
		dummyUser := entities.User{
			ID:       userId,
			Email:    "profile-user@example.com",
			Password: "hashedpassword123",
		}

		// Masukkan user dummy terlebih dahulu karena foreign key constraint
		err := database.DB.Create(&dummyUser).Error
		assert.NoError(t, err)

		dummyProfile := entities.Profile{
			ID:      "id-test-profile-001",
			UserId:  userId,
			Name:    "Irfan Aufa",
			Phone:   "081234567890",
			Address: "Jakarta",
		}

		err = profileRepo.DB.Create(&dummyProfile).Error
		assert.NoError(t, err)

		// Hapus data secara otomatis di akhir test (cleanup)
		// Karena relasi ON DELETE CASCADE, menghapus user akan otomatis menghapus profilnya
		t.Cleanup(func() {
			database.DB.Where("id = ?", userId).Delete(&entities.User{})
		})

		profile, err := profileRepo.FindProfileByUserId(userId)

		assert.NoError(t, err)
		assert.Equal(t, dummyProfile.ID, profile.ID)
		assert.Equal(t, dummyProfile.UserId, profile.UserId)
		assert.Equal(t, dummyProfile.Name, profile.Name)
		assert.Equal(t, dummyProfile.Phone, profile.Phone)
		assert.Equal(t, dummyProfile.Address, profile.Address)
	})
}

func Test_CreateProfile(t *testing.T) {
	t.Run("Create profile success", func(t *testing.T) {
		userId := "id-test-user-create-001"
		dummyUser := entities.User{
			ID:       userId,
			Email:    "create-profile-user@example.com",
			Password: "hashedpassword123",
		}

		// Buat user dummy untuk relasi foreign key
		err := database.DB.Create(&dummyUser).Error
		assert.NoError(t, err)

		t.Cleanup(func() {
			database.DB.Where("id = ?", userId).Delete(&entities.User{})
		})

		newProfile := &entities.Profile{
			ID:      "id-test-create-profile-001",
			UserId:  userId,
			Name:    "Fadilla",
			Phone:   "089876543210",
			Address: "Bandung",
		}

		err = profileRepo.CreateProfile(newProfile)
		assert.NoError(t, err)

		// Validasi tambahan: Pastikan profil benar-benar tersimpan di database
		var savedProfile entities.Profile
		err = profileRepo.DB.Where("user_id = ?", userId).First(&savedProfile).Error
		assert.NoError(t, err)
		assert.Equal(t, newProfile.ID, savedProfile.ID)
		assert.Equal(t, newProfile.Name, savedProfile.Name)
		assert.Equal(t, newProfile.Phone, savedProfile.Phone)
		assert.Equal(t, newProfile.Address, savedProfile.Address)
	})
}
