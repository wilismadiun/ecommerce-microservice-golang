package productclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ProductClientHttp struct{}

func (h *ProductClientHttp) DecreaseStock(productID string, quantity int) error {
	productServiceURL := os.Getenv("PRODUCTSERVICEURL")
	firstPath := "product"
	lastPath := "decrease-stock"

	targetURL := fmt.Sprintf("%s/%s/%s/%s", productServiceURL, firstPath, productID, lastPath)

	// Buat payload request sesuai kebutuhan controller tujuan
	payload := map[string]interface{}{
		"quantity": quantity,
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Gunakan http.NewRequest dengan method PATCH
	req, err := http.NewRequest(http.MethodPatch, targetURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Tangani jika response status code bukan 200 OK (misal 400 Bad Request)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to decrease stock, status code: %d", resp.StatusCode)
	}

	return nil
}
