package productclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestProductClientHttp_DecreaseStock(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return error when stock is less than quantity", func(t *testing.T) {
		// Mock server mensimulasikan error dari service tujuan (misal: Bad Request / Status 400)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPatch, r.Method)

			// Validasi payload yang masuk
			var body map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&body)

			// Simulasi jika quantity melebihi stok (service mengembalikan 400)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message":"insufficient stock"}`))
		}))
		defer mockServer.Close()

		// Set environment variable agar mengarah ke mock server
		t.Setenv("PRODUCTSERVICEURL", mockServer.URL)

		client := &ProductClientHttp{}
		err := client.DecreaseStock("prod-001", 10)

		// Assert bahwa client mendeteksi kegagalan/error
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		// Mock server mensimulasikan berhasil (Status 200 OK)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPatch, r.Method)

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"succes"}`))
		}))
		defer mockServer.Close()

		// Set environment variable
		t.Setenv("PRODUCTSERVICEURL", mockServer.URL)

		client := &ProductClientHttp{}
		err := client.DecreaseStock("prod-001", 3)

		// Assert bahwa proses berhasil tanpa error
		assert.NoError(t, err)
	})
}
