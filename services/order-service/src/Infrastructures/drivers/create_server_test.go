package drivers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"order-service/config"
	"order-service/src/Domains/entities"
	"order-service/src/Infrastructures/database"
	"order-service/src/Infrastructures/security"
	"os"
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

func Test_CreateOrder(t *testing.T) {
	// 1. Setup Mock Product Server untuk manipulasi DecreaseStock
	var productMockStatusCode = http.StatusOK
	mockProductServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(productMockStatusCode)
		if productMockStatusCode == http.StatusOK {
			_, _ = w.Write([]byte(`{"message":"success"}`))
		} else {
			_, _ = w.Write([]byte(`{"message":"insufficient stock"}`))
		}
	}))
	defer mockProductServer.Close()

	// Arahkan environment variable PRODUCTSERVICEURL ke mock server
	t.Setenv("PRODUCTSERVICEURL", mockProductServer.URL)

	// 2. Setup Router Gin
	router := gin.New()
	Router(router, database.DB)

	// 3. Generate Token JWT yang valid untuk pengujian
	tokenGenerator := security.AuthenticationTokenJWT{}
	validUserID := "user-test-uuid-123"
	validToken, err := tokenGenerator.GenerateToken(validUserID)
	assert.NoError(t, err)

	t.Run("should response 401 when authorization header is empty", func(t *testing.T) {
		payload := entities.Order{
			StoreID:         "store-test-999",
			ShippingAddress: "Jl. Merdeka No. 45",
			Items: []entities.OrderItem{
				{
					ProductID:   "prod-001",
					ProductName: "Mechanical Keyboard",
					Quantity:    2,
					Price:       1250000.0,
				},
			},
		}
		jsonBody, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/orders/"+validUserID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		// Tidak menyertakan Header Authorization

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("should response 400 when quantity empty", func(t *testing.T) {
		payload := entities.Order{
			StoreID:         "store-test-999",
			ShippingAddress: "Jl. Merdeka No. 45",
			Items: []entities.OrderItem{
				{
					ProductID:   "prod-001",
					ProductName: "Mechanical Keyboard",
					Quantity:    0, // Quantity invalid / kosong
					Price:       1250000.0,
				},
			},
		}
		jsonBody, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/orders/"+validUserID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+validToken)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("should response 400 when stock is less than quantity", func(t *testing.T) {
		// Manipulasi mock product server merespons 400 (stok tidak mencukupi)
		productMockStatusCode = http.StatusBadRequest

		payload := entities.Order{
			StoreID:         "store-test-999",
			ShippingAddress: "Jl. Merdeka No. 45",
			Items: []entities.OrderItem{
				{
					ProductID:   "prod-001",
					ProductName: "Mechanical Keyboard",
					Quantity:    10,
					Price:       1250000.0,
				},
			},
		}
		jsonBody, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/orders/"+validUserID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+validToken)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("should response 201 when order success", func(t *testing.T) {
		// Manipulasi mock product server merespons 200 (sukses mengurangi stok)
		productMockStatusCode = http.StatusOK

		payload := entities.Order{
			StoreID:         "store-test-999",
			ShippingAddress: "Jl. Merdeka No. 45",
			Items: []entities.OrderItem{
				{
					ProductID:   "prod-001",
					ProductName: "Mechanical Keyboard",
					Quantity:    2,
					Price:       1250000.0,
				},
			},
		}
		jsonBody, _ := json.Marshal(payload)

		// Bersihkan database setelah test selesai menggunakan t.Cleanup
		t.Cleanup(func() {
			database.DB.Unscoped().Where("user_id = ?", validUserID).Delete(&entities.Order{})
		})

		req, _ := http.NewRequest(http.MethodPost, "/orders/"+validUserID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+validToken)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)

		// Verifikasi data tersimpan di database
		var savedOrder entities.Order
		err := database.DB.Preload("Items").Where("user_id = ?", validUserID).First(&savedOrder).Error
		assert.NoError(t, err)
		assert.Equal(t, "PENDING", savedOrder.Status)
		assert.Equal(t, 2500000.0, savedOrder.TotalAmount)
		assert.Len(t, savedOrder.Items, 1)
	})
}
