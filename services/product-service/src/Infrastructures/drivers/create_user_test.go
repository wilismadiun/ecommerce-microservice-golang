package drivers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"product-service/config"
	"product-service/src/Domains/entities"
	"product-service/src/Infrastructures/database"
	"product-service/src/Infrastructures/security"
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

func Test_CreateProduct(t *testing.T) {
	// 1. Setup Mock HTTP Server Dinamis untuk StoreClient
	// Kita buat server yang status dan body-nya bisa diubah secara fleksibel per test case
	var currentStatusCode = http.StatusOK
	var currentResponse = "OWNER"

	mockStoreServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(currentStatusCode)
		if currentStatusCode == http.StatusOK {
			_, _ = w.Write([]byte(currentResponse))
		}
	}))
	defer mockStoreServer.Close()

	// Arahkan environment variable STORESERVICEURL ke mock server
	t.Setenv("STORESERVICEURL", mockStoreServer.URL)

	// 2. Setup Router Gin
	router := gin.New()
	Router(router, database.DB)

	// 3. Generate token JWT yang valid untuk pengujian
	tokenGenerator := security.AuthenticationTokenJWT{}
	validUserID := "user-test-uuid-123"
	validToken, err := tokenGenerator.GenerateToken(validUserID)
	assert.NoError(t, err)

	// Base payload produk yang valid (menggunakan entities.JSONMap untuk Attributes)
	validPayload := entities.Product{
		StoreID:     "store-test-999",
		Name:        "Mechanical Keyboard",
		Description: "Wireless RGB Mechanical Keyboard",
		Attributes: entities.JSONMap{
			"brand":  "Keychron",
			"switch": "Brown",
		},
		Price: 1250000.0,
		Stock: 25,
	}

	t.Run("should response 401 when authorization header is empty", func(t *testing.T) {
		jsonBody, _ := json.Marshal(validPayload)
		req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		// Tanpa Header Authorization

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("should response 404 when user role not found", func(t *testing.T) {
		currentStatusCode = http.StatusNotFound

		jsonBody, _ := json.Marshal(validPayload)
		req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+validToken)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})

	t.Run("should response 403 when user role not OWNER and ADMIN", func(t *testing.T) {
		currentStatusCode = http.StatusOK
		currentResponse = "CASHIER" // Role selain OWNER dan ADMIN

		jsonBody, _ := json.Marshal(validPayload)
		req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+validToken)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusForbidden, recorder.Code)
	})

	t.Run("should response 400 when name product empty", func(t *testing.T) {
		currentStatusCode = http.StatusOK
		currentResponse = "OWNER"

		// Buat salinan struct dengan nama kosong
		invalidPayload := validPayload
		invalidPayload.Name = ""

		jsonBody, _ := json.Marshal(invalidPayload)
		req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+validToken)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("create product success", func(t *testing.T) {
		currentStatusCode = http.StatusOK
		currentResponse = "OWNER"

		jsonBody, _ := json.Marshal(validPayload)
		req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+validToken)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		// Ambil ID produk yang baru saja dibuat dari response untuk dibersihkan (Cleanup)
		var response map[string]interface{}
		_ = json.Unmarshal(recorder.Body.Bytes(), &response)

		// Cleanup database otomatis setelah tes sukses berjalan
		t.Cleanup(func() {
			if data, ok := response["data"].(map[string]interface{}); ok {
				if prodID, exists := data["ID"].(string); exists {
					database.DB.Unscoped().Delete(&entities.Product{}, "id = ?", prodID)
				}
			}
		})
	})
}
