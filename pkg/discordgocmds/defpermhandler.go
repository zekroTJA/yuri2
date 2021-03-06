package discordgocmds

import (
	"github.com/zekroTJA/discordgo"
)

// DefaultPermissionHandler is the default handler
// for level based permission handling
type DefaultPermissionHandler struct {
	db DatabaseMiddleware
}

// NewDefaultPermissionHandler creates an instance of DefaultPermissionHandler
// getting passed the instance of the database middileware
func NewDefaultPermissionHandler(db DatabaseMiddleware) *DefaultPermissionHandler {
	return &DefaultPermissionHandler{
		db: db,
	}
}

// CheckUserPermission compares the command executing users permission level to the
// required permission level of the command and returns if the user matches the
// required permission.
func (p *DefaultPermissionHandler) CheckUserPermission(cmdArgs *CommandArgs, s *discordgo.Session, cmdInstance Command) (bool, error) {
	var roles []string
	var err error
	var lvl int

	if cmdArgs.User.ID == cmdArgs.CmdHandler.options.BotOwnerID {
		lvl = 999
	} else if cmdArgs.User.ID == cmdArgs.Guild.OwnerID {
		lvl = cmdArgs.CmdHandler.options.OwnerPermissionLevel
	} else {
		if cmdArgs.Guild.MemberCount > 100 {
			roles, err = getMemberRolesByRequest(s, cmdArgs.Guild.ID, cmdArgs.User.ID)
		} else {
			var found bool
			for _, m := range cmdArgs.Guild.Members {
				if m.User.ID == cmdArgs.User.ID {
					roles = m.Roles
					found = true
					break
				}
			}
			if !found {
				roles, err = getMemberRolesByRequest(s, cmdArgs.Guild.ID, cmdArgs.User.ID)
			}
		}

		if err != nil {
			return false, err
		}

		lvl, err = p.db.GetUserPermissionLevel(cmdArgs.User.ID, roles)
		if err != nil {
			return false, err
		}
	}

	return (cmdInstance.GetPermission() <= lvl), nil
}

// getMemberRolesByRequest tries to get the member by guildMember request and
// returns the members roles if successful.
func getMemberRolesByRequest(s *discordgo.Session, guildID, userID string) ([]string, error) {
	memb, err := s.GuildMember(guildID, userID)
	if err != nil {
		return nil, err
	}

	return memb.Roles, nil
}
