package storeclient

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"product-service/src/Domains/entities"
)

type StoreClientHttp struct{}

func (h *StoreClientHttp) GetUserRole(token, storeId string) (string, error) {
	storeServiceURL := os.Getenv("STORESERVICEURL")
	firstPath := "stores"
	lastPath := "user_role"

	targetURL := fmt.Sprintf("%s/%s/%s/%s", storeServiceURL, firstPath, storeId, lastPath)

	resp, err := http.Get(targetURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Memeriksa status code response (misal: jika bukan 200 OK)
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return "", entities.ErrUserRoleNotFound
		}

		return "", fmt.Errorf("failed to get user role, status code: %d", resp.StatusCode)
	}

	// Membaca isi body dari response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Mengubah body menjadi string (asumsi response dari service berupa teks/role langsung)
	role := string(bodyBytes)

	return role, nil
}
