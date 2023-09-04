package tokenx

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestJWTokenOk(t *testing.T) {
	r := require.New(t)
	secretKey := "testmainsecret"

	jwtTokenManager, err := NewJWTImpl(secretKey)
	r.NoError(err)

	username := randomManager.RandomString(8)
	issuedAt := time.Now()
	duration := time.Minute
	expiredAt := issuedAt.Add(duration)

	token, err := jwtTokenManager.CreateToken(username, duration)
	r.NoError(err)
	r.NotEmpty(token)

	payload, err := jwtTokenManager.VerifyToken(token)
	r.NoError(err)
	r.NotEmpty(payload)

	r.NotEmpty(payload.ID)
	r.Equal(username, payload.Username)
	r.WithinDuration(expiredAt, payload.ExpiredAt, time.Second)
	r.WithinDuration(issuedAt, payload.IssuedAt, time.Second)
}

func TestJWTokenExpired(t *testing.T) {
	r := require.New(t)
	secretKey := "testmainsecret"

	jwtTokenManager, err := NewJWTImpl(secretKey)
	r.NoError(err)

	username := randomManager.RandomString(8)
	duration := time.Minute

	token, err := jwtTokenManager.CreateToken(username, -duration)
	r.NoError(err)
	r.NotEmpty(token)

	payload, err := jwtTokenManager.VerifyToken(token)
	r.Error(err)
	r.EqualError(err, ErrTokenExpired.Error())

	r.Nil(payload)
}

func TestJWTokenNoHeaderVerify(t *testing.T) {
	r := require.New(t)

	username := randomManager.RandomString(8)
	duration := time.Minute
	secretKey := "testmainsecret"

	payload, err := NewPayload(username, duration)
	r.NoError(err)

	token := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	tokenString, err := token.SignedString(secretKey)
	r.NoError(err)
		
	jwtTokenManager, err := NewJWTImpl(secretKey)
	r.NoError(err)

	payload, err = jwtTokenManager.VerifyToken(tokenString)
	r.Error(err)
	r.EqualError(err, ErrInvalidToken.Error())

	r.Nil(payload)
}
