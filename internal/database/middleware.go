package database

import (
	"strings"
	"time"
)

// Timestamp is a uint8 array typed
// timetsmap which can be transformed
// to time.Time
type Timestamp []uint8

// ToTime parses the timestamp to a time object
func (t Timestamp) ToTime(format string) (time.Time, error) {
	tobj, err := time.Parse(format, string(t))
	if err != nil && strings.Contains(err.Error(), "out of range") {
		return time.Time{}, nil
	}
	return tobj, err
}

// SoundLogEntry is the data entity stored
// in the database as log entry when a sound
// was played.
type SoundLogEntry struct {
	Time    time.Time `json:"time"`
	UserID  string    `json:"user_id"`
	UserTag string    `json:"user_tag"`
	GuildID string    `json:"guild_id"`
	Source  string    `json:"source"`
	Sound   string    `json:"sound"`
}

type SoundStatsEntry struct {
	Sound string `json:"sound"`
	Count int    `json:"count"`
}

// Middleware describes the structure of a
// database middleware.
type Middleware interface {
	// Connect to the database server or file or
	// whatever you are about to use.
	Connect(params ...interface{}) error
	// Close the connection to the database.
	Close()
	GetConfigStructure() interface{}

	////////////////////////////
	// CMD HANDLER MIDDLEWARE //
	////////////////////////////

	// GetUserPermissionLevel returns the individual
	// permission level by the users ID and/or the
	// users role IDs.
	GetUserPermissionLevel(userID string, roles []string) (int, error)
	// GetGuildPrefix returns the individual prefix
	// for a guild by its ID.
	GetGuildPrefix(guildID string) (string, error)

	/////////////////
	// GUILD STUFF //
	/////////////////

	// SetGuildPrefix sets the custom prefix
	// for a guild in the DB.
	SetGuildPrefix(guildID, prefix string) error

	////////////////
	// USER STUFF //
	////////////////

	// SetFastTrigger gets the sound value used
	// for fast trigger. If this is an empty
	// string, this must be interpreted as
	// 'random sound'.
	GetFastTrigger(userID string) (string, error)
	// SetFastTrigger sets the sound which will
	// be triggered by using fast trigger.
	SetFastTrigger(userID, val string) error

	//////////////////
	// SOUNDS STUFF //
	//////////////////

	// AllLogEntry appends the log list by the
	// passed log data.
	AddLogEntry(sle *SoundLogEntry) error
	// GetLogEntries returns the log entries in
	// between the passed bounds. This list must
	// be ordered descending by time.
	GetLogEntries(guildID string, from, limit int) ([]*SoundLogEntry, error)
	// GetLogLen returns the ammount of entries
	// in the log wether per guildID, if passed
	// or of all entries.
	GetLogLen(guildID string) (int, error)

	// AddSoundStatsCount increases the play counter
	// of the sound for the specified guildID by one.
	AddSoundStatsCount(guildID, sound string) error
	// GetSoundStats returns the stats ordered
	// descending by play count.
	GetSoundStats(guildID string, limit int) ([]*SoundStatsEntry, error)

	SetGuildVolume(guildID string, volume int) error
	GetGuildVolume(guildID string) (int, error)
}
