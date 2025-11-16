package utils

import (
	"clean-archi/app/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken membuat JWT berdasarkan data alumni
func GenerateToken(alumni model.Alumni) (string, error) {
	claims := model.JWTClaims{
		UserID: alumni.MongoID.Hex(), // gunakan ObjectID Hex dari MongoDB
		Email:  alumni.Email,
		Role:   alumni.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "unair-auth-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-this" // fallback default
	}
	return token.SignedString([]byte(secret))
}

// ValidateToken memvalidasi token JWT dan mengembalikan klaim
func ValidateToken(tokenString string) (*model.JWTClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-this"
	}

	token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}
