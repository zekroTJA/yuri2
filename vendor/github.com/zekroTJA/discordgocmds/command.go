package discordgocmds

// Command is the interface containing
// Functions a command should have
type Command interface {
	GetInvokes() []string
	GetDescription() string
	GetHelp() string
	GetPermission() int
	Exec(args *CommandArgs) error
}
