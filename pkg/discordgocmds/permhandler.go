package discordgocmds

import (
	"github.com/zekroTJA/discordgo"
)

// PermissionHandler describes a struct handling
// permission of users using specific commands.
// By implementing this interface, you can create
// your own permission handler.
type PermissionHandler interface {
	// CheckUserPermission is getting passed the command arguments and the
	// instance of the command struct and returns a bool if the user is
	// permitted to use this command or an error, if the permission check
	// has failed for any reason, which will automatically count as
	// 'no permission'.
	CheckUserPermission(cmdArgs *CommandArgs, s *discordgo.Session, cmdInstance Command) (bool, error)
}
