package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zekroTJA/yuri2/internal/logger"
	"github.com/zekroTJA/yuri2/internal/static"

	"github.com/zekroTJA/yuri2/internal/database/sqlite"

	"github.com/ghodss/yaml"
	"github.com/zekroTJA/yuri2/internal/inits"
)

var (
	flagConfig     = flag.String("c", "./config.yml", "config file location")
	flagAddr       = flag.String("addr", "", "API expose address (overrides config)")
	flagDbDsn      = flag.String("db-dsn", "", "Database DSN (overrides config)")
	flagLLAddr     = flag.String("lavalink-address", "", "Lavalink address (overrides config)")
	flagLLPW       = flag.String("lavalink-password", "", "Lavalink password (overrides config)")
	flagLLSoundLoc = flag.String("lavalink-location", "", "Lavalink sounds location (overrides onfig)")
)

func main() {
	flag.Parse()

	/// CUSTOMIZABLE DRIVER SECTION ///////////////////////////

	// unmarshaler function for config
	unmarshaler := yaml.Unmarshal
	// marshaler function for preset config generation
	marshaler := func(v interface{}, prefix, indent string) ([]byte, error) {
		return yaml.Marshal(v)
	}

	// database middleware
	dbMiddleware := new(sqlite.SQLite)

	///////////////////////////////////////////////////////////

	// initializing teardown channel which will receive a
	// signal when one of the listed signals was sent to the
	// process or the program wants to exit itself by sending
	// a custom signal into the channel.
	teardownChan := make(chan os.Signal, 1)

	// init Logger
	inits.InitLogger()

	// init Config
	cfg := inits.InitConfig(*flagConfig, unmarshaler, marshaler, dbMiddleware.GetConfigStructure())

	if *flagAddr != "" {
		cfg.API.Address = *flagAddr
	}

	if *flagDbDsn != "" {
		cfg.Database = map[string]interface{}{
			"dsn": *flagDbDsn,
		}
	}

	if *flagLLAddr != "" {
		cfg.Lavalink.Address = *flagLLAddr
	}

	if *flagLLPW != "" {
		cfg.Lavalink.Password = *flagLLPW
	}

	if *flagLLSoundLoc != "" {
		cfg.Lavalink.SoundsLocations = []string{*flagLLSoundLoc}
	}

	// init Databse
	inits.InitDatabase(dbMiddleware, cfg.Database)
	// close database on exit
	defer func() {
		logger.Info("DATABASE :: shutting down")
		dbMiddleware.Close()
	}()

	// init Player
	player := inits.InitPlayer(cfg, dbMiddleware)

	// init Bot
	bot := inits.InitDiscordBot(cfg.Discord, dbMiddleware, player)
	// close bot connection on exit
	defer func() {
		logger.Info("DBOT :: shutting down")
		bot.Close()
	}()

	// init API
	api := inits.InitAPI(cfg, dbMiddleware, bot.Session, player, teardownChan)
	// close api exposure on exit
	defer func() {
		logger.Info("API :: shutting down")
		api.Close()
	}()

	// Block main thread until channel receives a teardown
	// signal. If the signal equals the custom signal
	// SigRestart which identifies a restart signal send
	// by Yuri hinself, the routine will be blocked for
	// 2 further seconds to ensure sending the response
	// to the request properly.
	signal.Notify(teardownChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	if sig := <-teardownChan; sig == static.SigRestart {
		logger.Info("CORE :: blocking core routine for 2 seconds to ensure restart request response")
		time.Sleep(2 * time.Second)
	}
}
