package storeclient

import (
	"net/http"
	"net/http/httptest"
	"product-service/src/Domains/entities"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StoreClient(t *testing.T) {
	token := "dummy-token"
	storeId := "store-123"

	t.Run("User role Not Found (404)", func(t *testing.T) {
		// 1. Buat mock HTTP server yang mengembalikan status 404
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		// 2. Set environment variable mengarah ke mock server URL
		t.Setenv("STORESERVICEURL", server.URL)

		client := &StoreClientHttp{}
		role, err := client.GetUserRole(token, storeId)

		// 3. Lakukan assertion
		assert.Error(t, err)
		assert.Equal(t, entities.ErrUserRoleNotFound, err)
		assert.Empty(t, role)
	})

	t.Run("Unexpected Status Code (500)", func(t *testing.T) {
		// 1. Buat mock HTTP server yang mengembalikan status 500
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		t.Setenv("STORESERVICEURL", server.URL)

		client := &StoreClientHttp{}
		role, err := client.GetUserRole(token, storeId)

		// 2. Lakukan assertion
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status code: 500")
		assert.Empty(t, role)
	})

	t.Run("user role exist (200)", func(t *testing.T) {
		expectedRole := "OWNER"

		// 1. Buat mock HTTP server yang mengembalikan status 200 dan body role
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(expectedRole))
		}))
		defer server.Close()

		t.Setenv("STORESERVICEURL", server.URL)

		client := &StoreClientHttp{}
		role, err := client.GetUserRole(token, storeId)

		// 2. Lakukan assertion
		assert.NoError(t, err)
		assert.Equal(t, expectedRole, role)
	})
}
