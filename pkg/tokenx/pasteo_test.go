package tokenx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasteoTokenOk(t *testing.T) {
	r := require.New(t)
	secretKey := randomManager.RandomString(32)

	pasteoManager, err := NewPasetoImpl(secretKey)
	r.NoError(err)

	username := randomManager.RandomString(8)
	issuedAt := time.Now()
	duration := time.Minute
	expiredAt := issuedAt.Add(duration)

	token, err := pasteoManager.CreateToken(username, duration)
	r.NoError(err)
	r.NotEmpty(token)

	payload, err := pasteoManager.VerifyToken(token)
	r.NoError(err)
	r.NotEmpty(payload)

	r.NotEmpty(payload.ID)
	r.Equal(username, payload.Username)
	r.WithinDuration(expiredAt, payload.ExpiredAt, time.Second)
	r.WithinDuration(issuedAt, payload.IssuedAt, time.Second)
}
