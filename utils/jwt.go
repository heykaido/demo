package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("exxo-secret")

func GenerateJWT(userID, name, email, roleType, customerID string) (string, error) {
	claims := jwt.MapClaims{
		"id":          userID,
		"name":        name,
		"email":       email,
		"role_type":   roleType,
		"customer_id": customerID,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
