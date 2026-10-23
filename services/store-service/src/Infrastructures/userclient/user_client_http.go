package userclient

import (
	"fmt"
	"net/http"
	"os"
)

type UserCLientHttp struct{}

func (h *UserCLientHttp) IsUserExists(userId string) (bool, error) {
	userServiceURL := os.Getenv("USERSERVICEURL")
	path := "profile"

	targetURL := fmt.Sprintf("%s/%s/%s", userServiceURL, userId, path)

	resp, err := http.Get(targetURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
