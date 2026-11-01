package productclient

type ProductClient interface {
	DecreaseStock(productID string, quantity int) error
}
