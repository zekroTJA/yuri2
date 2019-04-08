package auth

import (
	"time"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/pkg/random"
)

var tokenChars = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_.-")

type Auth struct {
	db database.Middleware

	hashingRounds int
	tokenLifetime time.Duration
}

func NewAuth(db database.Middleware, hashingRounds int, tokenLifetime time.Duration) *Auth {
	return &Auth{
		db:            db,
		hashingRounds: hashingRounds,
		tokenLifetime: tokenLifetime,
	}
}

func (a *Auth) CreateToken(userID string) (string, time.Time, error) {
	token := random.GetRandString(64, tokenChars)
	expires := time.Now().Add(a.tokenLifetime)
	hash, err := HashString(token, a.hashingRounds)
	if err != nil {
		return "", time.Time{}, err
	}

	if err = a.db.SetAuthToken(userID, hash, expires); err != nil {
		return "", time.Time{}, err
	}

	return token, expires, nil
}

func (a *Auth) CheckToken(userID, token string) (bool, error) {
	hash, err := a.db.GetAuthToken(userID)
	if err != nil {
		return false, err
	}

	return CompareHashString(hash, token), nil
}

func (a *Auth) RefreshToken(userID string) (time.Time, error) {
	expires := time.Now().Add(a.tokenLifetime)
	if err := a.db.SetAuthToken(userID, "", expires); err != nil {
		return time.Time{}, err
	}

	return expires, nil
}

func (a *Auth) CheckAndRefersh(userID, token string) (bool, time.Time, error) {
	ok, err := a.CheckToken(userID, token)
	if err != nil {
		return false, time.Time{}, err
	}
	if !ok {
		return false, time.Time{}, nil
	}

	expires, err := a.RefreshToken(userID)
	if err != nil {
		return false, time.Time{}, err
	}

	return true, expires, nil
}
