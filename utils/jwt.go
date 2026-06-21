package utils

import (
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID uint, email string) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(os.Getenv("JWT_SECRET")),
	)
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}