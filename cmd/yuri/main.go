package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/zekroTJA/yuri2/internal/logger"

	"github.com/zekroTJA/yuri2/internal/database/sqlite"

	"github.com/ghodss/yaml"
	"github.com/zekroTJA/yuri2/internal/inits"
)

var (
	flagConfig = flag.String("c", "./config.yml", "config file location")
)

func main() {
	flag.Parse()

	// unmarshaler function for config
	unmarshaler := yaml.Unmarshal
	// marshaler function for preset config generation
	marshaler := func(v interface{}, prefix, indent string) ([]byte, error) {
		return yaml.Marshal(v)
	}

	// database middleware
	dbMiddleware := new(sqlite.SQLite)

	// init Logger
	inits.InitLogger()

	// init Config
	cfg := inits.InitConfig(*flagConfig, unmarshaler, marshaler, dbMiddleware.GetConfigStructure())

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
	api := inits.InitAPI(cfg.API, dbMiddleware, bot.Session, player)
	// close api exposure on exit
	defer func() {
		logger.Info("API :: shutting down")
		api.Close()
	}()

	// block main go routine until one of the control
	// signals below was catched
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
