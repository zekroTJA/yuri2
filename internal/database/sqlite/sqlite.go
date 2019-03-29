package sqlite

import (
	"database/sql"
	"errors"

	"github.com/zekroTJA/yuri2/pkg/miltierror"

	// Importing sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// SQLite maintains the connection
// to the sqlite database file.
type SQLite struct {
	db *sql.DB
}

// Config contains the configuration
// values for the SQLite database
// connection.
type Config struct {
	DSN string `json:"dsn"`
}

// configFromMap creates a Config object
// from a map.
func configFromMap(m map[string]interface{}) (*Config, error) {
	var ok bool
	c := new(Config)

	if c.DSN, ok = m["dsn"].(string); !ok {
		return nil, errors.New("invalid config value type")
	}

	return c, nil
}

// Connect opens the database file and
// sets up database structure.
func (s *SQLite) Connect(params ...interface{}) error {
	var err error

	if len(params) < 1 {
		return errors.New("database file location musst be passed as first argument")
	}
	cfgMap, ok := params[0].(map[string]interface{})
	if !ok || cfgMap == nil {
		return errors.New("invalid parameter type or value")
	}

	cfg, err := configFromMap(cfgMap)
	if err != nil {
		return err
	}

	s.db, err = sql.Open("sqlite3", cfg.DSN)
	if err != nil {
		return err
	}

	return s.setup()
}

// setup sets up the database structure.
func (s *SQLite) setup() error {
	mErr := multierror.NewMultiError(nil)

	// TABLE `stats`
	_, err := s.db.Exec("CREATE TABLE IF NOT EXISTS `stats` (" +
		"`name` text NOT NULL DEFAULT ''," +
		"`guild_id` text NOT NULL DEFAULT ''," +
		"`played` int NOT NULL DEFAULT '0' );")
	mErr.Append(err)

	// TABLE `log`
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS `log` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT," +
		"`name` text NOT NULL DEFAULT ''," +
		"`time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, " +
		"`user_id` text NOT NULL DEFAULT '', " +
		"`guild_id` text NOT NULL DEFAULT '' );")
	mErr.Append(err)

	// TABLE `guilds`
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS `guilds` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT," +
		"`guild_id` text NOT NULL DEFAULT ''," +
		"`prefix` text NOT NULL DEFAULT '' );")
	mErr.Append(err)

	return mErr.Concat()
}

// Close the connection to the database.
func (s *SQLite) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

// GetConfigStructure returns an example object
// of the configuration structure for setting
// up the SQLite connection
func (s *SQLite) GetConfigStructure() interface{} {
	return &Config{
		DSN: "file:yuri.db.sqlite3",
	}
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
	var prefix string
	row := s.db.QueryRow("SELECT `prefix` FROM `guilds` WHERE `guild_id` = ?;", guildID)
	err := row.Scan(&prefix)
	if err == sql.ErrNoRows {
		return "", nil
	}

	return prefix, err
}

// SetGuildPrefix sets a prefix for a specific guild.
func (s *SQLite) SetGuildPrefix(guildID, prefix string) error {
	res, err := s.db.Exec("UPDATE `guilds` SET `prefix` = ? WHERE `guild_id` = ?;", prefix, guildID)
	if err != nil {
		return err
	}

	ar, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ar == 0 {
		_, err = s.db.Exec("INSERT INTO `guilds` (`guild_id`, `prefix`) VALUES (?, ?);", guildID, prefix)
	}

	return err
}
