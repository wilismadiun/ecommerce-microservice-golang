package security

type Hasher interface {
	Hash(password string) (string, error)
	CompareHashPassword(password, hashPassword string) error
}
