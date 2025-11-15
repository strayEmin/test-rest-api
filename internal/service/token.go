package service

import (
	"strconv"
	"test-task-rest-api/internal/config"
	"test-task-rest-api/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	config *config.Config
}

func NewTokenService(config *config.Config) *TokenService {
	return &TokenService{config: config}
}

type AccessClaims struct {
	jwt.RegisteredClaims
}

func (ts *TokenService) GenerateToken(user domain.User) (string, error) {
	now := time.Now().UTC()
	claims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    ts.config.Issuer,
			Subject:   strconv.FormatInt(user.ID, 10),
			ExpiresAt: jwt.NewNumericDate(now.Add(ts.config.Auth.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(ts.config.SecretKey)
}
