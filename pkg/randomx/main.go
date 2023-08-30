package randomx

import (
	"math/rand"
	"strings"
	"time"
)

var (
	alphabet     = "abcdefghijklmnopqrstuvwxyz"
)

type Random struct {
	*rand.Rand
}

func NewRandom() Random {
	return Random{
		rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Generate random int between min and max
func (r *Random) RandomInt(min, max int64) int64 {
	return min + r.Int63n(max-min+1)
}

// Generate random string of length n
func (r *Random) RandomString(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(alphabet[r.Intn(len(alphabet))])
	}
	return sb.String()
}