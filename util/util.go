package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// generates random numbers between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// generates random string of n characters
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// generates a random owner
func RandomOwner() string {
	return RandomString(6)
}

// generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// generates a random currency
func RandomCurrency() string {
	currency := []string{EUR, USD, CAD}
	n := len(currency)
	return currency[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomOwner())
}
