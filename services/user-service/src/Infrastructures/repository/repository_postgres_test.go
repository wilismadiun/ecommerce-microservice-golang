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

	code := m.Run()

	os.Exit(code)
}

func Test_FinsUserByEmail(t *testing.T) {
	t.Run("Email not found", func(t *testing.T) {
		email := "nonexistent@example.com"

		user, err := userRepo.FindUserByEmail(email)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrNotFound, err)
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
