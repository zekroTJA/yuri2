package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/zekroTJA/yuri2/internal/database"
	"github.com/zekroTJA/yuri2/pkg/miltierror"

	// Importing sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

const timeFormat = "2006-01-02 15:04:05"

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
		"`volume` int NOT NULL DEFAULT '100'," +
		"`prefix` text NOT NULL DEFAULT '' );")
	mErr.Append(err)

	// TABLE `users`
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS `users` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT," +
		"`user_id` text NOT NULL DEFAULT ''," +
		"`fast_trigger` text NOT NULL DEFAULT '' );")
	mErr.Append(err)

	// TABLE `sounds_log`
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS `sounds_log` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT," +
		"`time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"`user_id` text NOT NULL DEFAULT ''," +
		"`user_tag` text NOT NULL DEFAULT ''," +
		"`guild_id` text NOT NULL DEFAULT ''," +
		"`source` text NOT NULL DEFAULT ''," +
		"`sound` text NOT NULL DEFAULT '' );")
	mErr.Append(err)

	// TABLE `sounds_stats`
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS `sounds_stats` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT," +
		"`guild_id` text NOT NULL DEFAULT ''," +
		"`sound` text NOT NULL DEFAULT ''," +
		"`count` int NOT NULL DEFAULT '1' );")
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

/////////////////
// GUILD STUFF //
/////////////////

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

////////////////
// USER STUFF //
////////////////

func (s *SQLite) GetFastTrigger(userID string) (string, error) {
	var val string
	err := s.db.QueryRow("SELECT `fast_trigger` FROM `users` WHERE `user_id` = ?;", userID).Scan(&val)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return val, err
}

func (s *SQLite) SetFastTrigger(userID, val string) error {
	res, err := s.db.Exec("UPDATE `users` SET `fast_trigger` = ? WHERE `user_id` = ?;", val, userID)
	if err != nil {
		return err
	}

	ar, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ar == 0 {
		_, err = s.db.Exec("INSERT INTO `users` (`user_id`, `fast_trigger`) VALUES (?, ?);", userID, val)
	}

	return err
}

//////////////////
// SOUNDS STUFF //
//////////////////

func (s *SQLite) AddLogEntry(sle *database.SoundLogEntry) error {
	_, err := s.db.Exec("INSERT INTO `sounds_log` (`user_id`, `user_tag`, `guild_id`, `source`, `sound`) "+
		"VALUES (?, ?, ?, ?, ?);", sle.UserID, sle.UserTag, sle.GuildID, sle.Source, sle.Sound)
	return err
}

func (s *SQLite) GetLogEntries(guildID string, from, limit int) ([]*database.SoundLogEntry, error) {
	rows, err := s.db.Query("SELECT `time`, `user_id`, `user_tag`, `guild_id`, `source`, `sound` FROM `sounds_log` "+
		"WHERE `guild_id` = ? "+
		"ORDER BY `time` DESC LIMIT ?, ?;", guildID, from, limit)

	if err != nil {
		return nil, err
	}

	entries := make([]*database.SoundLogEntry, limit)
	i := 0
	for rows.Next() {
		sle := new(database.SoundLogEntry)

		err = rows.Scan(&sle.Time, &sle.UserID, &sle.UserTag, &sle.GuildID, &sle.Source, &sle.Sound)
		if err != nil {
			return nil, err
		}

		entries[i] = sle
		i++
	}

	if i < limit {
		return entries[:i], nil
	}

	return entries, nil
}

func (s *SQLite) GetLogLen(guildID string) (int, error) {
	var count int

	var row *sql.Row
	if guildID != "" {
		row = s.db.QueryRow("SELECT count(*) FROM `sounds_log` "+
			"WHERE `guild_id` = ?;", guildID)
	} else {
		row = s.db.QueryRow("SELECT count(*) FROM `sounds_log`;")
	}

	err := row.Scan(&count)

	return count, err
}

func (s *SQLite) AddSoundStatsCount(guildID, sound string) error {
	res, err := s.db.Exec("UPDATE `sounds_stats` SET `count` = `count` + 1 "+
		"WHERE `guild_id` = ? AND `sound` = ?;", guildID, sound)
	if err != nil {
		return err
	}

	ar, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ar == 0 {
		_, err = s.db.Exec("INSERT INTO `sounds_stats` (`guild_id`, `sound`) "+
			"VALUES (?, ?);", guildID, sound)
	}

	return err
}

func (s *SQLite) GetSoundStats(guildID string, limit int) ([]*database.SoundStatsEntry, error) {
	rows, err := s.db.Query("SELECT `sound`, `count` FROM `sounds_stats` "+
		"WHERE `guild_id` = ? ORDER BY `count` DESC "+
		"LIMIT ?;", guildID, limit)
	if err == sql.ErrNoRows {
		return make([]*database.SoundStatsEntry, 0), nil
	}
	if err != nil {
		return nil, err
	}

	res := make([]*database.SoundStatsEntry, limit)
	i := 0
	for rows.Next() {
		stat := new(database.SoundStatsEntry)
		err = rows.Scan(&stat.Sound, &stat.Count)
		if err != nil {
			return nil, err
		}
		res[i] = stat
		i++
	}

	fmt.Println(i)

	if i < limit {
		return res[:i], nil
	}

	return res, nil
}

func (s *SQLite) SetGuildVolume(guildID string, volume int) error {
	res, err := s.db.Exec("UPDATE `guilds` SET `volume` = ? WHERE `guild_id` = ?;", volume, guildID)
	if err != nil {
		return err
	}

	ar, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ar == 0 {
		_, err = s.db.Exec("INSERT INTO `guilds` (`guild_id`, `volume`) VALUES (?, ?);", guildID, volume)
	}

	return err
}

func (s *SQLite) GetGuildVolume(guildID string) (int, error) {
	var volume int
	row := s.db.QueryRow("SELECT `volume` FROM `guilds` WHERE `guild_id` = ?;", guildID)
	err := row.Scan(&volume)
	if err == sql.ErrNoRows {
		return 100, nil
	}

	return volume, err
}
