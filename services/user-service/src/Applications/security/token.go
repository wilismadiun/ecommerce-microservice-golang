package security

type AuthToken interface {
	GenerateToken(id string) (string, error)
}
