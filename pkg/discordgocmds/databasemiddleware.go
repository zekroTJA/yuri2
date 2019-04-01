package discordgocmds

// DatabaseMiddleware is an interface defining how
// to build a struct used by the command handler to
// communicate with your database, for example a
// MySql database.
type DatabaseMiddleware interface {
	// GetUserPermissionLevel returns the requested users
	// permission level number from the database or
	// and error, if the request failed for some reason.
	GetUserPermissionLevel(userID string, roles []string) (int, error)
	// GetGuildPrefix returns the requested guilds custom
	// prefix, if set. If the prefix was not set on the
	// guild, the function must return an empty string ("").
	// If the reuest failed for some reason, the function
	// returns an error.
	GetGuildPrefix(guildID string) (string, error)
}
