package apis

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/vietthangc1/simple_bank/pkg/envx"
	"github.com/vietthangc1/simple_bank/pkg/passwordx"
	"github.com/vietthangc1/simple_bank/pkg/randomx"
)

var (
	bycryptCost     = envx.Int("BYCRYPT_COST", 10)
	randomEntity    randomx.Random
	passwordManager passwordx.Password
)

func TestMain(m *testing.M) {
	gin.SetMode("release")

	randomEntity = randomx.NewRandom()
	passwordManager = passwordx.NewPassword(int(bycryptCost))

	os.Exit(m.Run())
}
