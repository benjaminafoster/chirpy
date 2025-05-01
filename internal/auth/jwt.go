package auth

import (
	"fmt"
	"net/http"
	"time"
	"strings"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy-access"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	issue_time := jwt.NewNumericDate(time.Now().UTC())
	expire_time := jwt.NewNumericDate(time.Now().UTC().Add(expiresIn))
	claims := jwt.RegisteredClaims{
		Issuer: string(TokenTypeAccess),
		IssuedAt: issue_time,
		ExpiresAt: expire_time,
		Subject: userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signingKey := []byte(tokenSecret)
	return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString, 
		&claimsStruct, 
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)

	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}

	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, fmt.Errorf("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return id, nil
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