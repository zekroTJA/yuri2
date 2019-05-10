package config

import (
	"io/ioutil"
	"os"
)

// UnmarshalFunc is the type of a function which
// can be used to parse byte data to an object
// instance.
type UnmarshalFunc func(data []byte, v interface{}) error

// MarshalIndentFunc is the type of a function which
// can be used to parse an object instance to raw
// data of a specific format in an indentet style,
// if wanted and supported by the function.
type MarshalIndentFunc func(v interface{}, prefix, indent string) ([]byte, error)

// Main contains the values of
// the parsed config file.
type Main struct {
	Discord  *Discord    `json:"discord"`
	Lavalink *Lavalink   `json:"lavalink"`
	Database interface{} `json:"database"`
	API      *API        `json:"api"`
	Misc     *Misc       `json:"misc"`
}

// Discord contains the values of
// Discrord specific configuration.
type Discord struct {
	Token          string            `json:"token"`
	OwnerID        string            `json:"owner_id"`
	GeneralPrefix  string            `json:"general_prefix"`
	RightRoleNames *DiscordRoleNames `json:"right_role_names"`
	StatusShuffle  *StatusShuffle    `json:"status_shuffle"`
}

// DiscordRoleNames contains the values of
// the permissioin role names for the
// Discord configuration.
type DiscordRoleNames struct {
	Player  string `json:"player"`
	Blocked string `json:"blocked"`
}

// Lavalink contains the values of
// the Lavalink specific configuration.
type Lavalink struct {
	Address         string   `json:"address"`
	Password        string   `json:"password"`
	SoundsLocations []string `json:"sounds_locations"`
}

// Misc contains miscellaneous configuration
// values.
type Misc struct {
	LogLevel int `json:"log_level"`
}

// StatusShuffle contains the configuration
// if the status shuffle.
type StatusShuffle struct {
	Status []string `json:"status"`
	Delay  string   `json:"delay"`
}

// API contains configuration for the
// REST and WS API.
type API struct {
	Enable        bool     `json:"enable"`
	ClientID      string   `json:"client_id"`
	ClientSecret  string   `json:"client_secret"`
	Address       string   `json:"address"`
	PublicAddress string   `json:"public_address"`
	TLS           *APITLS  `json:"tls"`
	AdminIDs      []string `json:"admin_ids"`
}

// APITLS contains configuration for the
// API TLS encryption.
type APITLS struct {
	Enable   bool   `json:"enable"`
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

// OpenAndParse tries to open the passed config file
// and parses the files content with the passed unmarshal
// function. If the config does not exist, a new one will
// be created with predefined default values.
//
// Parameters:
//   loc          : file location of the config
//   unamrshaler  : UnmarshalFunc to parse files content to an object
//   marshaler    : MarshalIndentFunc to parse the default config object
//                  to a new config file
//   dbConfStruct : database configuration object
//
// Returns:
//   *Main : the parsed config object - can be nil if parsing or opening fails
//   bool  : is true if config file was not existent and a new one was created
//   error : error if something fails
func OpenAndParse(loc string, unmarshaler UnmarshalFunc, marshaler MarshalIndentFunc, dbConfStruct interface{}) (*Main, bool, error) {
	data, err := ioutil.ReadFile(loc)
	if os.IsNotExist(err) {
		return nil, true, createNew(loc, marshaler, dbConfStruct)
	}
	if err != nil {
		return nil, false, err
	}

	c := new(Main)
	err = unmarshaler(data, c)
	return c, false, err
}

// createNEw creates a new config file with predefined default
// values parsed with the passed marshaler function.
//
// Parameters:
//   loc          : the config file location
//   marshaler    : MarshalIndentFunc to parse the default config object
//                  to a new config file
//   dbConfStruct : database configuration object
func createNew(loc string, marshaler MarshalIndentFunc, dbConfStruct interface{}) error {
	f, err := os.Create(loc)
	if err != nil {
		return err
	}

	var defMain = &Main{
		Discord: &Discord{
			GeneralPrefix: "y!",
			StatusShuffle: &StatusShuffle{
				Delay: "10s",
				Status: []string{
					"Yuri v.2!",
					"zekro.de",
					"github.com/zekroTJA/yuri2",
				},
			},
			RightRoleNames: &DiscordRoleNames{
				Player:  "@everyone",
				Blocked: "YuriBlocked",
			},
		},
		Lavalink: &Lavalink{
			Address:         "localhost:2333",
			SoundsLocations: make([]string, 0),
		},
		Database: dbConfStruct,
		API: &API{
			Enable:        false,
			Address:       ":443",
			PublicAddress: "https://yuri.example.com",
			AdminIDs:      make([]string, 0),
			TLS: &APITLS{
				Enable:   true,
				CertFile: "/etc/cert/example.com/example.com.cer",
				KeyFile:  "/etc/cert/example.com/example.com.key",
			},
		},
		Misc: &Misc{
			LogLevel: 4,
		},
	}

	data, err := marshaler(defMain, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	return err
}
