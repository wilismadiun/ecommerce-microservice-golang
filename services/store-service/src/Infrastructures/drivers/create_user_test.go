package drivers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"store-service/config"
	"store-service/src/Infrastructures/database"
	"store-service/src/Infrastructures/security"
	"testing"

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

func Test_AddStore(t *testing.T) {
	// 1. Setup Mock HTTP Server untuk User Service
	mockUserServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Jika path mengandung "user-not-found", kembalikan 404
		if r.URL.Path == "/user-not-found/profile" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// Untuk user lainnya, anggap user ditemukan (200 OK)
		w.WriteHeader(http.StatusOK)
	}))
	defer mockUserServer.Close()

	// Arahkan env USERSERVICEURL ke mock server
	os.Setenv("USERSERVICEURL", mockUserServer.URL)

	router := gin.New()
	Router(router, database.DB)

	tokenGenerator := &security.AuthenticationTokenJWT{}

	t.Run("should response 401 when authorization header is empty", func(t *testing.T) {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{"name": "My Store"})
		req, _ := http.NewRequest("POST", "/stores", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		// Sengaja tidak menyet header Authorization

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should response 404 when user not found", func(t *testing.T) {
		userId := "user-not-found"
		token, _ := tokenGenerator.GenerateToken(userId)

		w := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{"name": "My Store"})
		req, _ := http.NewRequest("POST", "/stores", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		// Menggunakan token/ID yang memicu respons 404 di mock server
		req.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should response 400 when store body invalid (empty name)", func(t *testing.T) {
		userId := "user-test-id-003"
		token, _ := tokenGenerator.GenerateToken(userId)

		w := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{"name": ""}) // Nama kosong melanggar VerifyStore
		req, _ := http.NewRequest("POST", "/stores", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should response 409 when store already exist", func(t *testing.T) {
		userId := "user-test-id-003"
		token, _ := tokenGenerator.GenerateToken(userId)

		// Setup: Masukkan data dummy store ke database
		database.DB.Exec(
			"INSERT INTO stores (id, name, description, phone, email, address) VALUES (?, ?, ?, ?, ?, ?)",
			"store-dup-1",
			"Existing Store",
			"Dummy Desc",
			"081234567890",
			"dummy@email.com",
			"Dummy Address",
		)
		defer database.DB.Exec("DELETE FROM stores WHERE id = 'store-dup-1'")

		w := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{
			"name":        "Existing Store",
			"description": "Deskripsi toko yang valid",
			"phone":       "081234567890",      // Panjang 10-15 karakter
			"email":       "store@example.com", // Format email valid
			"address":     "Jl. Contoh Alamat No. 123",
		})
		req, _ := http.NewRequest("POST", "/stores", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("should response 400 when user has 3 stores", func(t *testing.T) {
		userId := "user-limit-3"
		token, _ := tokenGenerator.GenerateToken(userId)
		// Setup: Masukkan 3 toko dan hubungkan ke user tersebut di database
		database.DB.Exec(`
			INSERT INTO stores (id, name, description, phone, email, address) 
    		VALUES 
		    ('s1', 'Store A', 'Desc Store A', '081234567890', 'storea@gmail.com', 'Alamat A'),
		    ('s2', 'Store B', 'Desc Store B', '081234567891', 'storeb@gmail.com', 'Alamat B'),
		    ('s3', 'Store C', 'Desc Store C', '081234567892', 'storec@gmail.com', 'Alamat C')
		`)
		database.DB.Exec("INSERT INTO user_stores (id, store_id, user_id, role_id) VALUES ('us1', 's1', ?, '1'), ('us2', 's2', ?, '1'), ('us3', 's3', ?, '1')", userId, userId, userId)

		// Cleanup database setelah test
		defer func() {
			database.DB.Exec("DELETE FROM user_stores WHERE user_id = ?", userId)
			database.DB.Exec("DELETE FROM stores WHERE id IN ('s1', 's2', 's3', 's4')")
		}()

		w := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{
			"name":        "Existing Store",
			"description": "Deskripsi toko yang valid",
			"phone":       "081234567890",      // Panjang 10-15 karakter
			"email":       "store@example.com", // Format email valid
			"address":     "Jl. Contoh Alamat No. 123",
		})
		req, _ := http.NewRequest("POST", "/stores", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should response 201 when add store success", func(t *testing.T) {
		userId := "valid-success-user"
		token, _ := tokenGenerator.GenerateToken(userId)

		// Cleanup database setelah test
		defer func() {
			database.DB.Exec("DELETE FROM user_stores")
			database.DB.Exec("DELETE FROM stores")
		}()

		w := httptest.NewRecorder()
		body, _ := json.Marshal(map[string]string{
			"name":        "Existing Store",
			"description": "Deskripsi toko yang valid",
			"phone":       "081234567890",      // Panjang 10-15 karakter
			"email":       "store@example.com", // Format email valid
			"address":     "Jl. Contoh Alamat No. 123",
		})
		req, _ := http.NewRequest("POST", "/stores", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}
