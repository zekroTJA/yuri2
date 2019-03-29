package discordgocmds

// CmdHandlerOptions is used to pass
// general options to the CmdHandler.
type CmdHandlerOptions struct {
	Prefix               string
	BotOwnerID           string
	OwnerPermissionLevel int
	DefaultColor         int
	InvokeToLower        bool
	ParseMsgEdit         bool
	ReactToBots          bool
	DeleteCmdMessages    bool
}

// NewCmdHandlerOptions creates a new instance
// of CmdHandlerOptions with the default
// configuration settings.
func NewCmdHandlerOptions() *CmdHandlerOptions {
	return &CmdHandlerOptions{
		Prefix:               "-",
		BotOwnerID:           "221905671296253953",
		OwnerPermissionLevel: 10,
		DefaultColor:         0x039BE5,
		InvokeToLower:        true,
		ParseMsgEdit:         true,
		ReactToBots:          false,
		DeleteCmdMessages:    true,
	}
}
