package drivers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"user-service/config"
	"user-service/src/Domains/entities"
	"user-service/src/Infrastructures/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Println("Peringatan: Gagal memuat .env dari root, mencoba path alternatif...")
	}

	dsn := config.DSN()
	database.ConnectDB(dsn)

	gin.SetMode(gin.TestMode)

	code := m.Run()

	os.Exit(code)
}

func Test_RegisterUser(t *testing.T) {
	router := gin.New()
	Router(router, database.DB)

	t.Run("should return error when verify user fails", func(t *testing.T) {
		payload := map[string]string{
			"email":    "",
			"password": "validpassword123",
		}
		jsonValue, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions (sesuaikan status code dengan handler Anda, misal 400 Bad Request)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when password has less than 8 characters", func(t *testing.T) {
		// Skenario: Password kurang dari 8 karakter
		payload := map[string]string{
			"email":    "shortpass@example.com",
			"password": "123",
		}
		jsonValue, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when email already exists", func(t *testing.T) {
		emailExists := "existing@example.com"

		// Seed data dummy ke database terlebih dahulu agar email sudah ada
		dummyUser := entities.User{
			ID:       "id-existing-001",
			Email:    emailExists,
			Password: "hashedpassword123",
		}
		err := database.DB.Create(&dummyUser).Error
		assert.NoError(t, err)

		// Hapus data secara otomatis di akhir sub-test ini
		t.Cleanup(func() {
			database.DB.Where("email = ?", emailExists).Delete(&entities.User{})
		})

		payload := map[string]string{
			"email":    emailExists,
			"password": "validpassword123",
		}
		jsonValue, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Biasanya konflik (409) atau bad request (400) tergantung implementasi handler
		assert.True(t, w.Code == http.StatusBadRequest || w.Code == http.StatusConflict)
	})

	t.Run("create user success", func(t *testing.T) {
		emailSuccess := "success-register@example.com"

		// Hapus data secara otomatis di akhir sub-test (cleanup) jika berhasil masuk DB
		t.Cleanup(func() {
			database.DB.Where("email = ?", emailSuccess).Delete(&entities.User{})
		})

		payload := map[string]string{
			"email":    emailSuccess,
			"password": "validpassword123",
		}
		jsonValue, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions sukses (bisa 200 OK atau 201 Created tergantung handler Anda)
		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusCreated)

		// Validasi tambahan: Pastikan data benar-benar tersimpan di database
		var savedUser entities.User
		err := database.DB.Where("email = ?", emailSuccess).First(&savedUser).Error
		assert.NoError(t, err)
		assert.Equal(t, emailSuccess, savedUser.Email)
	})
}
