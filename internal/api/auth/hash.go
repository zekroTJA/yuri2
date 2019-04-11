package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashString creates a blowfish generated hash
// of the passed string using the passed number
// of rounds as cost.
func HashString(str string, rounds int) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(str), rounds)
	return string(b), err
}

// CompareHashString compares a hashed
// stirng with a clear string using the
// blowfish hash algorithm.
func CompareHashString(hash, str string) bool {
	res := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	return res == nil
}
