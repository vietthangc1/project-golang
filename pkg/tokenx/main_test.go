package tokenx

import (
	"testing"

	"github.com/vietthangc1/simple_bank/pkg/randomx"
)

var (
	randomManager randomx.Random
)

func TestMain(m *testing.M) {
	randomManager = randomx.NewRandom()
}
