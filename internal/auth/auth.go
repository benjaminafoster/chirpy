package auth

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

func HashPassword(password string) (string, error) {
	log.Printf("Password: %s", password)

	hashed_pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	
	log.Printf("Hashed password: %s", hashed_pwd)
	return string(hashed_pwd), nil
}

func CheckPasswordHash(storedHash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		return err
	}

	return nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no authorization header present")
	}

	tokenFields := strings.Fields(authHeader)
	if len(tokenFields) != 2 {
		return "", fmt.Errorf("authorization header must follow convention: 'Bearer <token>'")
	}

	return tokenFields[1], nil
}

