package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var signingKey []byte

func Init() {
	key, ok := os.LookupEnv("SIGNING_KEY")
	if !ok {
		panic("SIGNING_KEY does not exist in environment variables!")
	}

	signingKey = []byte(key)
}

func New(id string) (string, error) {
	claims := jwt.StandardClaims{
		Id:       id,
		IssuedAt: time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signingKey)
}

func Parse(tokenString string) (*string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("failed to convert to map claims")
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	id, ok := claims["jti"].(string)

	if !ok {
		return nil, fmt.Errorf("failed to convert jti")
	}

	return &id, nil
}
