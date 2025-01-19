package util

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateJWT(userID int, username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	secretKey := os.Getenv("SECRET_KEY_JWT")
	fmt.Println(secretKey)

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
