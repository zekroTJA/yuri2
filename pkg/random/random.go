package random

import (
	"crypto/rand"
	"math/big"
)

var chars = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// GetRandString returnes a random stirng
// of the defined length.
func GetRandString(n int, charSet []rune) (string, error) {
	if charSet == nil {
		charSet = chars
	}

	nb, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
	if err != nil {
		return "", err
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = charSet[nb.Int64()]
	}

	return string(b), nil
}
