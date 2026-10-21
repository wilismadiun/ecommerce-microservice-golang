package security

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = os.Getenv("JWT_SECRET")

type AuthenticationTokenJWT struct{}

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func (h *AuthenticationTokenJWT) GenerateToken(id, role string) (string, error) {
	now := time.Now()

	claims := CustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Hour)), // Jangan lupa aktifkan jika perlu
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   id, // 'id' biasanya diletakkan di Subject (sub)
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (h *AuthenticationTokenJWT) ValidateToken(tokenString string) (string, string, error) {

	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("invalid signing method")
			}

			return []byte(secretKey), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", "", errors.New("Token sudah tidak berlaku")
		}
		return "", "", err
	}

	if !token.Valid {
		return "", "", jwt.ErrTokenInvalidClaims
	}

	return claims.Subject, claims.Role, nil
}
