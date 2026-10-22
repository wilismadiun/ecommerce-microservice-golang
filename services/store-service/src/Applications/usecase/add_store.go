package usecase

import (
	"errors"
	"store-service/src/Applications/generator"
	"store-service/src/Applications/userclient"
	"store-service/src/Domains/entities"
	"store-service/src/Domains/repository"
)

type AddStore struct {
	StoreRepository     repository.StoreRepository
	UserStoreRepository repository.UserStoreRepository
	Generator           generator.GeneratorId
	UserCLient          userclient.UserClient
}

func (h *AddStore) Execute(userId string, store entities.Store) (entities.Store, error) {
	// verify user
	userExist, err := h.UserCLient.IsUserExists(userId)
	if err != nil {
		return entities.Store{}, err
	}

	if !userExist {
		return entities.Store{}, entities.ErrUserNotFound
	}

	// verify store
	err = entities.VerifyStore(store)
	if err != nil {
		return entities.Store{}, err
	}

	// check store name in database
	_, err = h.StoreRepository.FindStoreByName(store.Name)
	if err == nil {
		return entities.Store{}, entities.ErrStoreAlreadyExists
	}

	if !errors.Is(err, entities.ErrStoreNotFound) {
		return entities.Store{}, err
	}

	// Check how many stores the user owns
	storesOwned, err := h.UserStoreRepository.FindUserStoreByUser(userId)
	if err != nil {
		return entities.Store{}, err
	}

	if len(storesOwned) >= 3 {
		return entities.Store{}, errors.New("Max store limit reached")
	}

	// generate store id
	store.Id = h.Generator.Generator()

	// add store in database
	err = h.StoreRepository.AddStore(&store)
	if err != nil {
		return entities.Store{}, err
	}

	// add owner in database
	userStoreId := h.Generator.Generator()
	userStore := entities.UserStore{
		Id:      userStoreId,
		StoreId: store.Id,
		UserId:  userId,
		RoleId:  "1",
	}

	err = h.UserStoreRepository.AddUserStore(userStore)
	if err != nil {
		_ = h.StoreRepository.DeleteStoreById(store.Id)
		return entities.Store{}, err
	}

	return store, nil
}
