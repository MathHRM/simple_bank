package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
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
	currencies := []string{BRL, USD, KWZ}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomMoney() int64 {
	return RandomInt(1, 1000)
}

func RandomName() string {
	return RandomString(7)
}

func RandomEmail(name string) string {
	emailDomain := []string{"@email.com", "@gmail.net", "@quentemail.br", "@salazard.wizard"}
	n := len(emailDomain)
	return name + emailDomain[rand.Intn(n)]
}