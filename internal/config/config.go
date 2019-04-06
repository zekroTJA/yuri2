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
	Misc     *Misc       `json:"misc"`
}

// Discord contains the values of
// Discrord specific configuration.
type Discord struct {
	Token         string         `json:"token"`
	OwnerID       string         `json:"owner_id"`
	GeneralPrefix string         `json:"general_prefix"`
	StatusShuffle *StatusShuffle `json:"status_shuffle"`
}

// Lavalink contains the values of
// the Lavalink specific configuration.
type Lavalink struct {
	Address        string `json:"address"`
	Password       string `json:"password"`
	SoundsLocation string `json:"sounds_location"`
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
		},
		Lavalink: &Lavalink{
			Address: "localhost:2333",
		},
		Database: dbConfStruct,
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
