package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomString(n int) string {
	var stringBuilder strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[r.Intn(k)]
		stringBuilder.WriteByte(c)
	}

	return  stringBuilder.String()
}

func RandomInt(min, max int64) int64 {
	return min + r.Int63n(max - min + 1)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomBalance() int64 {
	return RandomInt(1, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "IDR"}
	k := len(currencies)

	return currencies[r.Intn(k)]
}
