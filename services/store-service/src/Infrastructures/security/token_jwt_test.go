package security

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AuthToken(t *testing.T) {
	id := "user-123"
	generatorToken := &AuthenticationTokenJWT{}

	authToken, err := generatorToken.GenerateToken(id)
	log.Println("ini adlah token")
	log.Println(authToken)

	assert.NoError(t, err)
	assert.NotEmpty(t, authToken)

	idClaim, err := generatorToken.ValidateToken(authToken)
	log.Println("ini adalah claim")
	log.Println(idClaim)

	assert.NoError(t, err)
	assert.Equal(t, id, idClaim)
}
