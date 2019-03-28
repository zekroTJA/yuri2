package sqlite

import (
	"database/sql"
	"errors"

	// Importing sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// SQLite maintains the connection
// to the sqlite database file.
type SQLite struct {
	db *sql.DB
}

// Connect opens the database file and
// sets up database structure.
func (s *SQLite) Connect(params ...interface{}) error {
	var err error

	if len(params) < 1 {
		return errors.New("database file location musst be passed as first argument")
	}
	dsn, ok := params[0].(string)
	if !ok || dsn == "" {
		return errors.New("invalid parameter type or value")
	}

	s.db, err = sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}

	return s.setup()
}

// setup sets up the database structure.
func (s *SQLite) setup() error {
	// TODO: FUNCTIONALITY
	return nil
}

// Close the connection to the database.
func (s *SQLite) Close() {
	s.db.Close()
}

////////////////////////////
// CMD HANDLER MIDDLEWARE //
////////////////////////////

// GetUserPermissionLevel returns the individual
// permission level by the users ID and/or the
// users role IDs.
func (s *SQLite) GetUserPermissionLevel(userID string, roles []string) (int, error) {
	// TODO: FUNCTIONALITY
	return 0, nil
}

// GetGuildPrefix returns the individual prefix
// for a guild by its ID.
func (s *SQLite) GetGuildPrefix(guildID string) (string, error) {
	// TODO: FUNCTIONALITY
	return "", nil
}
