package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AuthToken(t *testing.T) {
	id := "user-123"
	generatorToken := &AuthenticationTokenJWT{}

	authToken, err := generatorToken.GenerateToken(id)

	assert.NoError(t, err)
	assert.NotEmpty(t, authToken)

	idClaim, err := generatorToken.ValidateToken(authToken)

	assert.NoError(t, err)
	assert.Equal(t, id, idClaim)
}
