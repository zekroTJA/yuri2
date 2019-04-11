package auth

import (
	"time"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/pkg/random"
)

// token chars to be used for token generation
var tokenChars = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_.-")

// Auth provides methods to create authorization tokens,
// saving them as blowfish hashed string in the database
// checking them or refreshing expiring time of them.
type Auth struct {
	db database.Middleware

	hashingRounds int
	tokenLifetime time.Duration
}

// NewAuth creates a new instance of Auth.
func NewAuth(db database.Middleware, hashingRounds int, tokenLifetime time.Duration) *Auth {
	return &Auth{
		db:            db,
		hashingRounds: hashingRounds,
		tokenLifetime: tokenLifetime,
	}
}

// CreateToken randomly generates a token string,
// which will be hashed by the blowfish algorithm
// using the specified rounds and saves the hash
// in combination with the userID and the calcuated
// expire time (from specified expire duration)
// to the database.
// Return values are the the clear token, the expire
// time and an error object if something failed.
func (a *Auth) CreateToken(userID string) (string, time.Time, error) {
	token, err := random.GetRandString(64, tokenChars)
	if err != nil {
		return "", time.Time{}, err
	}

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

// CheckToken checks the passed token by getting the
// saved hash from the DB by userID and checking it
// against the passed token.
// This will return true, if the check was successful.
// AN non-nil error will only be returned if the
// database access fails.
func (a *Auth) CheckToken(userID, token string) (bool, error) {
	entry, err := a.db.GetAuthToken(userID)
	if err != nil || entry == nil {
		return false, err
	}

	return CompareHashString(entry.TokenHash, token), nil
}

// RefreshToken gets the current time and adds the specified
// expire duration to it. This will be set as expiring time
// for the token which is specified for the passed userID.
func (a *Auth) RefreshToken(userID string) (time.Time, error) {
	expires := time.Now().Add(a.tokenLifetime)
	if err := a.db.SetAuthToken(userID, "", expires); err != nil {
		return time.Time{}, err
	}

	return expires, nil
}

// CheckAndRefresh is shorthand for CheckToken and, if
// this check was passed successfully, refreshing the
// expire time of the token.
func (a *Auth) CheckAndRefresh(userID, token string) (bool, time.Time, error) {
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
