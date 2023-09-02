package passwordx

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vietthangc1/simple_bank/pkg/envx"
	"github.com/vietthangc1/simple_bank/pkg/randomx"
)

var (
	cost = envx.Int("BYCRYPT_COST", 10)
	passwordManager Password
	randomManager randomx.Random
)

func TestMain(m *testing.M) {
	passwordManager = NewPassword(int(cost))
	randomManager = randomx.NewRandom()
}

func TestPassword(t *testing.T) {
	r := require.New(t)
	randomPassword := randomManager.RandomString(8)

	hashedPassword, err := passwordManager.HashPassword(randomPassword)
	r.NoError(err)
	r.NotEmpty(hashedPassword)

	err = passwordManager.CheckPassword(randomPassword, hashedPassword)
	r.NoError(err)
}


