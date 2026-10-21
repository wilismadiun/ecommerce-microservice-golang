package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AuthToken(t *testing.T) {
	id := "user-123"
	role := "admin"
	generatorToken := &AuthenticationTokenJWT{}

	authToken, err := generatorToken.GenerateToken(id, role)

	assert.NoError(t, err)
	assert.NotEmpty(t, authToken)

	idClaim, roleClaim, err := generatorToken.ValidateToken(authToken)

	assert.NoError(t, err)
	assert.Equal(t, id, idClaim)
	assert.Equal(t, role, roleClaim)
}
