package util

import (
	"math/rand"
	"strings"

	"github.com/shopspring/decimal"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomFloat() float64 {

	numA := rand.Int31()
	numB := rand.Int31()

	return rand.Float64() * float64(numA-numB)
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

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "TWD"}

	n := len(currencies)

	return currencies[rand.Intn(n)]
}

func RandomBalance() decimal.Decimal {

	return decimal.NewFromFloat(RandomFloat()).Round(3)
}
