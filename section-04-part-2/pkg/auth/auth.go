package auth

import (
	"errors"
	"fmt"
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

func ParseToken(authenticationHeader string) (string, error) {
	tokenString := authenticationHeader[len("Bearer "):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("cannot parse")
		}
		secret := []byte("ct-secret-key")
		fmt.Println(secret)
		return secret, nil
	})
	fmt.Println(token)
	if err != nil {
		return "", errors.New(fmt.Sprintf("invalid token: %v", err))
	}
	if token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		username := claims["sub"].(string)
		return username, nil
	}
	return "", errors.New("invalid token")
}
