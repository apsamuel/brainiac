package main

import (
	"flag"

	"github.com/apsamuel/brainiac/pkg/cache"
	"github.com/apsamuel/brainiac/pkg/database"
	"github.com/apsamuel/brainiac/pkg/logger"
)

var ConfigFile string
var Debug bool

func init() {
	flag.StringVar(&ConfigFile, "config", "./config/brainiac.yaml", "Path to the configuration file")
	flag.BoolVar(&Debug, "debug", false, "Enable debug mode")
	flag.Parse()
}

/*
 * This is the main entry point for the brainiac application.
 */
func main() {

	l := logger.Logger
	l.Logger.Info().Msg("starting brainiac")

	var databaseConfig database.Config
	var cacheConfig cache.Config
	if Debug {
		l.Logger.Debug().Msg("debugging enabled")
	}

	err := databaseConfig.Configure(ConfigFile)
	if err != nil {
		l.Logger.Error().Msg(err.Error())
	}
	databaseConfig.Log = &l.Logger

	_, err = database.MakeStorage(databaseConfig)
	if err != nil {
		l.Logger.Error().Msg(err.Error())
	}

	err = cacheConfig.Configure(ConfigFile)
	if err != nil {
		l.Logger.Error().Msg(err.Error())
	}

	cacheConfig.Log = &l.Logger
	_, err = cache.MakeStorage(cacheConfig)
	if err != nil {
		l.Logger.Error().Msg(err.Error())
	}
}
