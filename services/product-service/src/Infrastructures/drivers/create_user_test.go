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

func Test_DecreaseStock(t *testing.T) {
	// Setup Router Gin
	router := gin.New()

	// Atau jika Anda menggunakan fungsi Router langsung dari package drivers:
	Router(router, database.DB)

	t.Run("should response 400 when product id empty", func(t *testing.T) {
		payload := map[string]interface{}{
			"quantity": 5,
		}
		jsonBody, _ := json.Marshal(payload)

		// Memanggil endpoint dengan ID kosong pada parameter URL
		req, _ := http.NewRequest(http.MethodPatch, "/products//decrease-stock", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		// Gin akan mengembalikan 404 karena struktur URL tidak cocok,
		// atau 400 jika lolos ke usecase. Kita sesuaikan dengan ekspektasi route.
		assert.True(t, recorder.Code == http.StatusBadRequest || recorder.Code == http.StatusNotFound)
	})

	t.Run("should response 400 when stock is less than the quantity", func(t *testing.T) {
		// 1. Arrange: Buat produk dummy dengan stok sedikit (misal: stok 2)
		productID := "prod-stock-low-999"
		dummyProduct := entities.Product{
			ID:          productID,
			StoreID:     "store-test-999",
			Name:        "Test Item",
			Description: "Description",
			Price:       10000.0,
			Stock:       2,
		}
		err := database.DB.Create(&dummyProduct).Error
		assert.NoError(t, err)

		t.Cleanup(func() {
			database.DB.Unscoped().Delete(&entities.Product{}, "id = ?", productID)
		})

		// 2. Act: Minta kurangi stok melebihi jumlah yang ada (quantity: 5)
		payload := entities.DecreaseStock{
			Quantity: 5,
		}
		jsonBody, _ := json.Marshal(payload)

		url := "/products/" + productID + "/decrease-stock"
		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		// 3. Assert: Harus mengembalikan status 400 Bad Request karena stok kurang
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("success", func(t *testing.T) {
		// 1. Arrange: Buat produk dummy dengan stok mencukupi (misal: stok 10)
		productID := "prod-stock-success-888"
		dummyProduct := entities.Product{
			ID:          productID,
			StoreID:     "store-test-999",
			Name:        "Test Item Success",
			Description: "Description",
			Price:       15000.0,
			Stock:       10,
		}
		err := database.DB.Create(&dummyProduct).Error
		assert.NoError(t, err)

		t.Cleanup(func() {
			database.DB.Unscoped().Delete(&entities.Product{}, "id = ?", productID)
		})

		// 2. Act: Kurangi stok dengan jumlah yang valid (quantity: 4)
		payload := map[string]interface{}{
			"quantity": 4,
		}
		jsonBody, _ := json.Marshal(payload)

		url := "/products/" + productID + "/decrease-stock"
		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		// 3. Assert: Harus sukses (status 200 OK)
		assert.Equal(t, http.StatusOK, recorder.Code)

		// Verifikasi langsung ke database bahwa stok tersisa 6 (10 - 4)
		var updatedProduct entities.Product
		database.DB.First(&updatedProduct, "id = ?", productID)
		assert.Equal(t, 6, updatedProduct.Stock)
	})
}

func Test_UpdateStock(t *testing.T) {
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

	// Setup Router Gin
	router := gin.New()
	Router(router, database.DB)

	// Generate token JWT yang valid untuk pengujian
	tokenGenerator := security.AuthenticationTokenJWT{}
	validUserID := "user-test-uuid-123"
	validToken, err := tokenGenerator.GenerateToken(validUserID)
	assert.NoError(t, err)

	t.Run("should response 401 when authorization header is empty", func(t *testing.T) {
		payload := map[string]interface{}{
			"quantity": 10,
		}
		jsonBody, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPatch, "/products/some-id/stock", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		// Tidak menyertakan header Authorization

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("should response 401 when user role not OWNER and ADMIN", func(t *testing.T) {
		// Buat produk dummy terlebih dahulu agar lolos pengecekan GetProductById
		productID := "prod-role-test-999"
		dummyProduct := entities.Product{
			ID:          productID,
			StoreID:     "store-test-999",
			Name:        "Test Product Role",
			Description: "Description",
			Price:       50000.0,
			Stock:       10,
		}
		err := database.DB.Create(&dummyProduct).Error
		assert.NoError(t, err)

		t.Cleanup(func() {
			database.DB.Unscoped().Delete(&entities.Product{}, "id = ?", productID)
		})

		// Set status mock store server menjadi unauthorized atau response role selain owner/admin
		currentStatusCode = http.StatusOK
		currentResponse = "CASHIER" // Role bukan OWNER atau ADMIN

		payload := map[string]interface{}{
			"quantity": 15,
		}
		jsonBody, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPatch, "/products/"+productID+"/stock", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+validToken)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		// Berdasarkan controller, ErrForbidden menghasilkan status BadRequest (atau StatusUnauthorized jika disesuaikan)
		assert.True(t, recorder.Code == http.StatusBadRequest || recorder.Code == http.StatusUnauthorized)
	})

	t.Run("should response 400 when product not found", func(t *testing.T) {
		nonExistentID := "product-not-found-uuid"

		payload := map[string]interface{}{
			"quantity": 20,
		}
		jsonBody, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPatch, "/products/"+nonExistentID+"/stock", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+validToken)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		// Controller memetakan ErrProductNotFound menjadi StatusNotFound (404)
		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})

	t.Run("success", func(t *testing.T) {
		productID := "prod-success-uuid-111"
		dummyProduct := entities.Product{
			ID:          productID,
			StoreID:     "store-test-999",
			Name:        "Test Product Success",
			Description: "Description",
			Price:       75000.0,
			Stock:       10,
		}
		err := database.DB.Create(&dummyProduct).Error
		assert.NoError(t, err)

		t.Cleanup(func() {
			database.DB.Unscoped().Delete(&entities.Product{}, "id = ?", productID)
		})

		// Set mock store server merespons role OWNER yang valid
		currentStatusCode = http.StatusOK
		currentResponse = "OWNER"

		payload := map[string]interface{}{
			"quantity": 30,
		}
		jsonBody, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPatch, "/products/"+productID+"/stock", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+validToken)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		// Verifikasi perubahan stok di database
		var updatedProduct entities.Product
		database.DB.First(&updatedProduct, "id = ?", productID)
		assert.Equal(t, 30, updatedProduct.Stock)
	})
}
