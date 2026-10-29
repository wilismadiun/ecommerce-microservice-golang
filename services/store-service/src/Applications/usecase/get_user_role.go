package usecase

import "store-service/src/Domains/repository"

type GetUserRole struct {
	UserStoreRepository repository.UserStoreRepository
}

func (h *GetUserRole) Execute(userId, storeId string) (string, error) {
	role, err := h.UserStoreRepository.GetUserRole(userId, storeId)
	if err != nil {
		return "", err
	}

	return role, nil
}
