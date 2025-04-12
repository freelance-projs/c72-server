package auth

import (
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ngoctd314/c72-api-server/pkg/model"
)

type Claims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

func NewClaims(opts ...ClaimOption) *Claims {
	now := time.Now()

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	for _, opt := range opts {
		opt(claims)
	}

	return claims
}

func generateToken(opts ...ClaimOption) (string, error) {
	claims := NewClaims(opts...)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKeyBytes, _ := base64.StdEncoding.DecodeString("xxx")
	accessToken, err := jwtToken.SignedString(secretKeyBytes)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func ParseToken(token string) (*Claims, error) {
	claims := &Claims{}

	secretKeyBytes, _ := base64.StdEncoding.DecodeString("xxx")
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return secretKeyBytes, nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}

type ClaimOption func(*Claims)

func withUsername(username string) ClaimOption {
	return func(c *Claims) {
		c.Username = username
	}
}

func withRole(role model.ERole) ClaimOption {
	return func(c *Claims) {
		c.Role = role.String()
	}
}
