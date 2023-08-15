package util

import (
	"math/rand"
	"strings"

	"github.com/shopspring/decimal"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomInt(min, max int64) int64{
	return min + rand.Int63n(max - min + 1)
}

func Randomstring(n int) string {
	var sb strings.Builder

	K := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(K)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return Randomstring(6)
}

func RandomBalance() decimal.Decimal {
	return decimal.NewFromInt(RandomInt(0, 1000))
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "TWD"}

	n := len(currencies)

	return currencies[rand.Intn(n)]
}