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

func NewUser(id string) (string, error) {
	claims := jwt.StandardClaims{
		Id:       id,
		IssuedAt: time.Now().Unix(),
		Subject:  "user",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signingKey)
}

func NewGuest(id string, expiresAt time.Time) (string, error) {
	claims := jwt.StandardClaims{
		Id:        id,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: expiresAt.Unix(),
		Subject:   "guest",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signingKey)
}

func parse(tokenString string) (jwt.MapClaims, error) {
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

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	if !ok {
		return nil, fmt.Errorf("failed to convert to map claims")
	}

	return claims, nil
}

func ParseUser(tokenString string) (*string, error) {
	claims, err := parse(tokenString)

	if err != nil {
		return nil, err
	}

	subject, ok := claims["sub"].(string)

	if !ok {
		return nil, fmt.Errorf("failed to convert sub")
	}

	if subject != "user" {
		return nil, fmt.Errorf("subject is not user")
	}

	id, ok := claims["jti"].(string)

	if !ok {
		return nil, fmt.Errorf("failed to convert jti")
	}

	return &id, nil
}

func ParseGuest(tokenString string) (*string, error) {
	claims, err := parse(tokenString)

	if err != nil {
		return nil, err
	}

	subject, ok := claims["sub"].(string)

	if !ok {
		return nil, fmt.Errorf("failed to convert sub")
	}

	if subject != "guest" {
		return nil, fmt.Errorf("subject is not guest")
	}

	fmt.Printf("exp: %d", claims["exp"])

	expiresAt, ok := claims["exp"].(float64)

	if !ok {
		return nil, fmt.Errorf("failed to convert exp")
	}

	if time.Now().Unix() > int64(expiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	id, ok := claims["jti"].(string)

	if !ok {
		return nil, fmt.Errorf("failed to convert jti")
	}

	return &id, nil
}
