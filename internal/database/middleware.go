package database

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
}
