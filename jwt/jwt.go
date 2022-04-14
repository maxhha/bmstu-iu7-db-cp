package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var signingKey []byte

var ErrWrongSigningMethod = errors.New("unexpected signing method")

func Init() {
	key, ok := os.LookupEnv("SIGNING_KEY")
	if !ok {
		panic("SIGNING_KEY does not exist in environment variables!")
	}

	signingKey = []byte(key)
}

func NewUser(id string) (string, error) {
	claims := jwt.StandardClaims{
		Id:       id,
		IssuedAt: time.Now().Unix(),
		Subject:  "user",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signingKey)
}

func parse(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrWrongSigningMethod // fmt.Errorf("%w: %v", ErrWrongSigningMethod, token.Header["alg"])
		}

		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to convert to map claims")
	}

	return claims, nil
}

func ParseUser(tokenString string) (string, error) {
	claims, err := parse(tokenString)
	if err != nil {
		return "", err
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("failed to convert sub")
	}

	if subject != "user" {
		return "", fmt.Errorf("subject is not user")
	}

	id, ok := claims["jti"].(string)
	if !ok {
		return "", fmt.Errorf("failed to convert jti")
	}

	return id, nil
}
