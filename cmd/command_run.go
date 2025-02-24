package cmd

import (
	"github.com/apsamuel/brainiac/pkg/cache"
	"github.com/apsamuel/brainiac/pkg/common"
	"github.com/apsamuel/brainiac/pkg/database"
	"github.com/apsamuel/brainiac/pkg/extensions/api"
	"github.com/apsamuel/brainiac/pkg/logger"
	"github.com/spf13/cobra"
)

var runCommand = &cobra.Command{
	Use:   "run",
	Short: "Run the brainiac server",
	Long:  `Run the brainiac server`,
	Run: func(cmd *cobra.Command, args []string) {
		var databaseConfig database.Config
		var cacheConfig cache.Config
		var apiConfig api.Config

		var brainiacConfig common.Config

		// debug brainiacConfig

		observerChannels := make(map[string]chan database.Item)
		for _, observer := range []string{"api", "proxy"} {
			observerChannels[observer] = make(chan database.Item)
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

		err = databaseConfig.ConfigureFromFile(configFile)
		if err != nil {
			l.Logger.Error().Msg(err.Error())
		}
		databaseConfig.Log = &l.Logger

		storage, err := database.MakeStorage(databaseConfig)
		if err != nil {
			l.Logger.Error().Msg(err.Error())
		}

		err = cacheConfig.ConfigureFromFile(configFile)
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

		apiHandler := api.Agent{
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
	},
}

func init() {
}
