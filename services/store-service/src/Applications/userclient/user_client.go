package userclient

type UserClient interface {
	IsUserExists(userId string) (bool, error)
}
