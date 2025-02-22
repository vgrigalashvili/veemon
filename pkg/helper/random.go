package helper

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	GEL = "GEL"
	USD = "USD"
	EUR = "EUR"
)
const alphabet = "abcdefghijklmnopqrstuvwxyz"

var rnd *rand.Rand

func init() {
	seed := time.Now().UnixNano()
	rnd = rand.New(rand.NewSource(seed))
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rnd.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rnd.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{GEL, USD, EUR}
	n := len(currencies)
	return currencies[rnd.Intn(n)]
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

// RandomPin generates a random 6-digit pin
// And returns it as a int
func RandomPin() int {
	return int(RandomInt(100000, 999999))
}
