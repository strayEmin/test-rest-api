package jwt

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type TokenMetadata struct {
	UserID    int64
	Issuer    string
	ExpiresAt int64
	IssuedAt  int64
	NotBefore int64
}

func ExtractTokenMetadata(c *fiber.Ctx) (TokenMetadata, error) {
	jwtCtx := c.Locals("jwt").(*jwt.Token)
	claims := jwtCtx.Claims.(jwt.MapClaims)

	UserID, err := strconv.ParseInt(claims["sub"].(string), 10, 64)
	if err != nil {
		return TokenMetadata{}, err
	}

	return TokenMetadata{
		UserID:    UserID,
		Issuer:    claims["iss"].(string),
		ExpiresAt: int64(claims["exp"].(float64)),
		IssuedAt:  int64(claims["iat"].(float64)),
		NotBefore: int64(claims["nbf"].(float64)),
	}, nil
}
