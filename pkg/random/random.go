package random

import (
	"math/rand"
	"time"
)

var chars = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetRandString returnes a random stirng
// of the defined length.
func GetRandString(n int, charSet []rune) string {
	if charSet == nil {
		charSet = chars
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}
