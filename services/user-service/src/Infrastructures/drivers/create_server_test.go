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
	"user-service/src/Infrastructures/security"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
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

func Test_Login(t *testing.T) {
	router := gin.New()
	Router(router, database.DB)

	t.Run("should response 400 when email is empty", func(t *testing.T) {
		payload := map[string]string{
			"email":    "",
			"password": "validpassword123",
		}
		jsonValue, _ := json.Marshal(payload)

		// Sesuaikan path "/login" dengan routing di aplikasi Anda jika berbeda
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should response 404 when email is not registered", func(t *testing.T) {
		payload := map[string]string{
			"email":    "unregistered@example.com",
			"password": "validpassword123",
		}
		jsonValue, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should response 401 when password is incorrect", func(t *testing.T) {
		email := "wrongpass@example.com"
		// Generate hash password yang valid di database untuk user ini
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword123"), bcrypt.DefaultCost)

		dummyUser := entities.User{
			ID:       "id-login-001",
			Email:    email,
			Password: string(hashedPassword),
		}

		// Masukkan data dummy ke database
		err := database.DB.Create(&dummyUser).Error
		assert.NoError(t, err)

		// Hapus data secara otomatis di akhir sub-test
		t.Cleanup(func() {
			database.DB.Where("email = ?", email).Delete(&entities.User{})
		})

		payload := map[string]string{
			"email":    email,
			"password": "wrongpassword", // Password salah yang dikirimkan
		}
		jsonValue, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("login success", func(t *testing.T) {
		email := "successlogin@example.com"
		plainPassword := "correctpassword123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

		dummyUser := entities.User{
			ID:       "id-login-002",
			Email:    email,
			Password: string(hashedPassword),
		}

		// Masukkan data dummy ke database
		err := database.DB.Create(&dummyUser).Error
		assert.NoError(t, err)

		// Hapus data secara otomatis di akhir sub-test
		t.Cleanup(func() {
			database.DB.Where("email = ?", email).Delete(&entities.User{})
		})

		payload := map[string]string{
			"email":    email,
			"password": plainPassword, // Password benar
		}
		jsonValue, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func Test_AddProfile(t *testing.T) {
	router := gin.New()
	Router(router, database.DB)

	// Helper untuk membuat token JWT valid berdasarkan User ID
	tokenGenerator := &security.AuthenticationTokenJWT{}

	t.Run("should response 401 when authorization header is empty", func(t *testing.T) {
		payload := map[string]string{
			"name":    "John Doe",
			"phone":   "081234567890",
			"address": "Jl. Merdeka No. 1",
		}
		jsonValue, _ := json.Marshal(payload)

		// Request tanpa header Authorization
		req, _ := http.NewRequest(http.MethodPost, "/profile", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should response 400 when phone number with fewer than 10 digits", func(t *testing.T) {
		userId := "user-test-id-002"
		token, _ := tokenGenerator.GenerateToken(userId)

		payload := map[string]string{
			"name":    "John Doe",
			"phone":   "12345", // Kurang dari 10 digit (gagal verifikasi profile)
			"address": "Jl. Merdeka No. 1",
		}
		jsonValue, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/profile", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should response 409 when the profile already exists", func(t *testing.T) {
		userId := "user-test-id-003"
		token, _ := tokenGenerator.GenerateToken(userId)

		dummyUser := entities.User{
			ID:       userId,
			Email:    "testuser003@example.com",
			Password: "hashedpassword",
		}
		err := database.DB.Create(&dummyUser).Error
		assert.NoError(t, err)

		// Seed (masukkan) data profile dummy terlebih dahulu ke database
		dummyProfile := entities.Profile{
			ID:      "profile-id-existing",
			UserId:  userId,
			Name:    "John Doe",
			Phone:   "081234567890",
			Address: "Jl. Merdeka No. 1",
		}
		err = database.DB.Create(&dummyProfile).Error
		assert.NoError(t, err)

		// Cleanup data profile setelah test selesai
		t.Cleanup(func() {
			database.DB.Where("user_id = ?", userId).Delete(&entities.Profile{})
			database.DB.Where("id = ?", userId).Delete(&entities.User{})
		})

		payload := map[string]string{
			"name":    "John Doe Updated",
			"phone":   "081234567890",
			"address": "Jl. Sudirman No. 2",
		}
		jsonValue, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/profile", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("should return 201 when profile successfully created", func(t *testing.T) {
		userId := "user-test-id-004"
		token, _ := tokenGenerator.GenerateToken(userId)

		dummyUser := entities.User{
			ID:       userId,
			Email:    "testuser003@example.com",
			Password: "hashedpassword",
		}
		err := database.DB.Create(&dummyUser).Error
		assert.NoError(t, err)

		// Cleanup data profile setelah test sukses dimasukkan ke database
		t.Cleanup(func() {
			database.DB.Where("user_id = ?", userId).Delete(&entities.Profile{})
			database.DB.Where("id = ?", userId).Delete(&entities.User{})
		})

		payload := map[string]string{
			"name":    "Jane Doe",
			"phone":   "089876543210",
			"address": "Jl. Asia Afrika No. 5",
		}
		jsonValue, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/profile", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Validasi tambahan: Cek langsung ke database apakah data benar-benar tersimpan
		var savedProfile entities.Profile
		err = database.DB.Where("user_id = ?", userId).First(&savedProfile).Error
		assert.NoError(t, err)
		assert.Equal(t, "Jane Doe", savedProfile.Name)
		assert.Equal(t, "089876543210", savedProfile.Phone)
	})
}

func Test_GetProfileByUser(t *testing.T) {
	router := gin.New()
	Router(router, database.DB)

	t.Run("should response 404 when profile not found", func(t *testing.T) {
		nonExistentUserID := "999999"

		req, _ := http.NewRequest(http.MethodGet, "/users/"+nonExistentUserID+"/profile", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("should response 200 when profile found", func(t *testing.T) {
		testUserID := "user-test-123"
		testProfileID := "prof-test-123"

		// --- 1. SETUP: Masukkan data ke 'users' terlebih dahulu (karena ada foreign key) ---
		database.DB.Exec(
			"INSERT INTO users (id, email, password) VALUES (?, ?, ?)",
			testUserID, "test@mail.com", "hashedpassword123",
		)

		// --- 2. SETUP: Masukkan data ke 'profiles' ---
		database.DB.Exec(
			"INSERT INTO profiles (id, user_id, name, phone, address) VALUES (?, ?, ?, ?, ?)",
			testProfileID, testUserID, "John Doe", "08123456789", "Jakarta Selatan",
		)

		t.Cleanup(func() {
			database.DB.Exec("DELETE FROM profiles WHERE user_id = ?", testUserID)
			database.DB.Exec("DELETE FROM users WHERE id = ?", testUserID)
		})

		req, _ := http.NewRequest(http.MethodGet, "/users/"+testUserID+"/profile", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Profile retrieved successfully")
	})
}
