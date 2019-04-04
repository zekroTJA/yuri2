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

	unmarshaler := yaml.Unmarshal
	marshaler := func(v interface{}, prefix, indent string) ([]byte, error) {
		return yaml.Marshal(v)
	}

	dbMiddleware := new(sqlite.SQLite)
	defer func() {
		logger.Info("DATABASE :: shutting down")
		dbMiddleware.Close()
	}()

	// init Logger
	inits.InitLogger()
	// init Config
	cfg := inits.InitConfig(*flagConfig, unmarshaler, marshaler, dbMiddleware.GetConfigStructure())
	// init Databse
	inits.InitDatabase(dbMiddleware, cfg.Database)
	// init Player
	player := inits.InitPlayer(cfg.Lavalink)
	// init Bot
	bot := inits.InitDiscordBot(cfg.Discord.Token, cfg.Discord.OwnerID,
		cfg.Discord.GeneralPrefix, dbMiddleware, player)
	defer func() {
		logger.Info("DBOT :: shutting down")
		bot.Close()
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
