package auth

import (
	"golang.org/x/crypto/bcrypt"
	"log"
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