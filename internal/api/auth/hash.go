package auth

import "golang.org/x/crypto/bcrypt"

func HashString(str string, rounds int) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(str), rounds)
	return string(b), err
}

func CompareHashString(hash, str string) bool {
	res := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	return res == nil
}
