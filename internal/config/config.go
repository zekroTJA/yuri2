package config

// Main contains the values of
// the parsed config file.
type Main struct {
	Discord *Discord `json:"discord"`
}

// Discord contains the values of
// Discrord specific configuration.
type Discord struct {
	Token         string `json:"token"`
	OwnerID       string `json:"owner_id"`
	GeneralPrefix string `json:"general_prefix"`
}
