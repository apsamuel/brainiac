package main

import (
	"flag"

	"github.com/apsamuel/brainiac/pkg/cache"
	"github.com/apsamuel/brainiac/pkg/common"
	"github.com/apsamuel/brainiac/pkg/database"
	"github.com/apsamuel/brainiac/pkg/extensions/api"
	"github.com/apsamuel/brainiac/pkg/logger"
)

var configFile string
var debug bool

func init() {
	flag.StringVar(&configFile, "config", "./config/brainiac.yaml", "Path to the configuration file")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.Parse()
}

/*
 * This is the main entry point for the brainiac application.
 */
func main() {
	var databaseConfig database.Config
	var cacheConfig cache.Config
	var apiConfig api.Config

	var brainiacConfig common.Config

	// debug brainiacConfig

	observerChannels := make(map[string]chan common.Item)
	for _, observer := range []string{"api", "proxy"} {
		observerChannels[observer] = make(chan common.Item)
	}
	l := logger.Logger
	l.Logger.Info().Msg("starting brainiac")

	if debug {
		l.Logger.Debug().Msg("debugging enabled")
	}

	_, err := brainiacConfig.FromFile(configFile)
	if err != nil {
		logger.Logger.Logger.Error().Msg(err.Error())
		return
	}
	l.Logger.Info().Interface("configuration", brainiacConfig).Msg("loaded brainiac configuration")
	jsonConfig, err := brainiacConfig.ToInterface()
	if err != nil {
		l.Logger.Error().Msg(err.Error())
		return
	}
	l.Logger.Info().Interface("json configuration", jsonConfig).Msg("loaded brainiac configuration")

	err = databaseConfig.Configure(configFile)
	if err != nil {
		l.Logger.Error().Msg(err.Error())
	}
	databaseConfig.Log = &l.Logger

	storage, err := database.MakeStorage(databaseConfig)
	if err != nil {
		l.Logger.Error().Msg(err.Error())
	}

	err = cacheConfig.Configure(configFile)
	if err != nil {
		l.Logger.Error().Msg(err.Error())
	}

	cacheConfig.Log = &l.Logger
	cacheStorage, err := cache.MakeStorage(cacheConfig)
	if err != nil {
		l.Logger.Error().Msg(err.Error())
	}

	/* configure and start API server */

	err = apiConfig.Configure(configFile)
	if err != nil {
		l.Logger.Error().Msg(err.Error())
	}

	apiConfig.Log = &l.Logger

	apiHandler := api.Handler{
		Config:    &apiConfig,
		Observers: observerChannels,
		Storage:   storage,
		Cache:     &cacheStorage,
	}
	err = apiHandler.MakeRouter()
	if err != nil {
		l.Logger.Error().Msg(err.Error())
	}

	for _, agentRoutes := range [][]*common.Route{
		apiHandler.ListRoutes(),
	} {
		for _, route := range agentRoutes {
			err := apiHandler.AddRoute(route)
			if err != nil {
				l.Logger.Error().Msg(err.Error())
				return
			}
		}
	}

	go func() {
		l.Logger.Info().Msg("starting brainiac api server")
		err := apiHandler.Serve()
		if err != nil {
			l.Logger.Error().Msg(err.Error())
			return
		}
	}()

	select {}
}
