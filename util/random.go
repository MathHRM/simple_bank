package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n( max - min + 1 )
}

func RandomString(leng int) string {
	var sb strings.Builder
	var k int = len(alphabet)

	for i := 0; i < leng; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomCurrency() string {
	currencies := []string{"USD", "BRL", "EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}