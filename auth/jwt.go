package auth

import (
	"RuntimeError/types/domain"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var signingKey = []byte("mega-super-secret-key")

type JWTClaims struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func newJwt(user *types.User) (string, error) {
	claims := &JWTClaims{
		UserId: user.Id,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 60 * time.Minute)),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(signingKey)
}

func parseJwt(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid jwt token")
	}

	return claims, nil
}

type JWTProvider struct{}

func (p *JWTProvider) Create(ctx context.Context, user *types.User) (string, error) {
	return newJwt(user)
}

func (p *JWTProvider) Verify(ctx context.Context, token string) (interface{}, error) {
	res, err := parseJwt(token)
	if err != nil {
		return nil, err
	}

	return res, nil
}
