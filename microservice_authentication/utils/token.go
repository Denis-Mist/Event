package utils

import (
	"auth/models"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func generateSecretKey() (string, error) {
	b := make([]byte, 32) // 32 bytes is a good size for a secret key
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func GenerateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	key, err := generateSecretKey()
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
