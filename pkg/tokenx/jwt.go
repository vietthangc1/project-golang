package tokenx

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/vietthangc1/simple_bank/pkg/envx"
)

var (
	minSecretKeyLength = envx.Int("MIN_JWT_SECRET_SIZE", 10)
)

type JWTImpl struct {
	secretKey string
}

func NewJWTImpl(secretKey string) (Token, error) {
	if len(secretKey) < int(minSecretKeyLength) {
		return nil, fmt.Errorf("secret key must be at least %d chracters", minSecretKeyLength)
	}
	return &JWTImpl{
		secretKey: secretKey,
	}, nil
}

func (j *JWTImpl) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString(j.secretKey)
}

func (j *JWTImpl) VerifyToken(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
