package token

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	GenerateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

type jwtService struct{}

func NewJWTService() JWTService {
	return &jwtService{}
}

func (service *jwtService) GenerateToken(username string, duration time.Duration) (string, error) {
	secret := os.Getenv("TOKEN_SYMMETRIC_KEY")
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(secret))
	return token, err
}

// VerifyToken checks if the token is valid or not
func (service *jwtService) VerifyToken(token string) (*Payload, error) {
	secret := os.Getenv("TOKEN_SYMMETRIC_KEY")
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
