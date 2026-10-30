package storeclient

type StoreClient interface {
	GetUserRole(token, storeId string) (string, error)
}
