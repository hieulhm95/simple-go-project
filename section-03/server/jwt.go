package main

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func GenerateToken(username string, expireDuration time.Duration) (string, error) {
	mySigningKey := []byte("ct-secret-key")

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "ct-backend-course",
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}
