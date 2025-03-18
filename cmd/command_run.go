package cmd

import (
	"fmt"
	"os"
	"strconv"

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

		l := logger.Logger
		l.Logger.Info().Msg("starting brainiac")

		if debug {
			l.Logger.Debug().Msg("debugging enabled")
		}

		var configPort int
		if configFile == "" {
			if configEngine == "postgres" {
				l.Logger.Info().Msgf("using config host %s", configHost)
				/*
					configHost defines the FQDN or IP of the configuration database server
					if the value is not provided the program will panic
					this value can be set by providing the --configHost flag or by setting the BRAINIAC_CONFIG_HOST environment variable
				*/
				if configHost == "" {
					configHost = os.Getenv("BRAINIAC_CONFIG_HOST")
					if configHost == "" {
						panic("BRAINIAC_CONFIG_HOST not set")
					}

				}
				/*
					configPort defines the port of the configuration database server
					if the value is not provided the program will panic
					this value can be set by providing the --configPort flag or by setting the BRAINIAC_CONFIG_PORT environment variable
				*/
				l.Logger.Info().Msgf("using config port %d", configPort)
				if configPort == 0 {
					configPortStr := os.Getenv("BRAINIAC_CONFIG_PORT")
					if configPortStr == "" {
						panic("BRAINIAC_CONFIG_PORT is not set")
					}
					// if configPort 1 - 1024 - privileged ports
					var portError error
					configPort, portError = strconv.Atoi(configPortStr)
					if portError != nil {
						panic(portError)
					}

				}

				/*
					configDatabase defines the name of the configuration database
					if the value is not provided the program will panic
					this value can be set by providing the --configDatabase flag or by setting the BRAINIAC_CONFIG_DB environment variable
					NOTE: if using a redis configuration backend, this value can be set to 0
				*/
				l.Logger.Info().Msgf("using config database %s", configDatabase)
				if configDatabase == "" {
					configDatabase = os.Getenv("BRAINIAC_CONFIG_DB")
					if configDatabase == "" {
						panic("BRAINIAC_CONFIG_DB not set")
					}
				}

				/*
					configTable defines the table or key for the configuration data
					if the value is not provided the program will panic
					this value can be set by providing the --configTable flag or by setting the BRAINIAC_CONFIG_TABLE environment variable
				*/
				l.Logger.Info().Msgf("using config table %s", configTable)
				if configTable == "" {
					configTable = os.Getenv("BRAINIAC_CONFIG_TABLE")
					if configTable == "" {
						panic("BRAINIAC_CONFIG_TABLE not set")
					}
				}

				/*
					configUsername defines the username for the configuration database
					if the value is not provided the program will panic
					this value is set by providing the BRAINIAC_CONFIG_USER environment variable
				*/

				configUsername := os.Getenv("BRAINIAC_CONFIG_USER")
				if configUsername == "" {
					panic("BRAINIAC_CONFIG_USER not set")
				}
				l.Logger.Info().Msgf("using config username %s", configUsername)
				/*
					configPassword defines the password for the configuration database
					if the value is not provided the program will panic
					this value is set by providing the BRAINIAC_CONFIG_PASS environment variable
				*/
				configPassword := os.Getenv("BRAINIAC_CONFIG_PASS")
				if configPassword == "" {
					panic("BRAINIAC_CONFIG_PASS not set")
				}

				/*
					configKey defines the key for the AES encryption
					if the value is not provided the program will panic
					this value is set by providing the BRAINIAC_AES_KEY environment variable

					this value should be a 32 byte key
				*/
				configKey := os.Getenv(aesKeyVariable)
				if configKey == "" {
					panic(fmt.Sprintf("%s not set", aesKeyVariable))
				}

				/*
					configNonce defines the nonce for the AES encryption
					if the value is not provided the program will panic
					this value is set by providing the BRAINIAC_AES_NONCE environment variable

					this value should be a 12 byte nonce
				*/
				configNonce := os.Getenv(aesNonceVariable)
				if configNonce == "" {
					panic(fmt.Sprintf("%s not set", aesNonceVariable))
				}

				config, err := database.RetrieveConfig(
					configHost,
					configPort,
					configDatabase,
					"config_data",
					configUsername,
					configPassword,
					configKey,
					configNonce,
				)
				if err != nil {
					panic(err)
				}
				_, err = brainiacConfig.FromBytes(config)
				if err != nil {
					l.Logger.Error().Msg(err.Error())
				}
				// l.Logger.Info().Interface("config", string(config)).Msg("retrieved configuration")
			}

		} else {
			l.Logger.Info().Str("config file", configFile).Msg("loading configuration file")
			_, err := brainiacConfig.FromFile(configFile)
			if err != nil {
				logger.Logger.Logger.Error().Msg(err.Error())
				return
			}
			l.Logger.Info().Interface("configuration", brainiacConfig).Msg("loaded brainiac configuration")
		}

		jsonConfig, err := brainiacConfig.ToInterface()

		observerChannels := database.GetObservers(jsonConfig)
		if err != nil {
			l.Logger.Error().Msg(err.Error())
			return
		}
		l.Logger.Info().Interface("json configuration", jsonConfig).Msg("loaded brainiac configuration")

		err = databaseConfig.ConfigureFromInterface(jsonConfig)
		if err != nil {
			l.Logger.Error().Msg(err.Error())
		}
		databaseConfig.Log = &l.Logger
		storage, err := database.MakeStorage(databaseConfig)
		if err != nil {
			l.Logger.Error().Msg(err.Error())
		}

		err = cacheConfig.ConfigureFromInterface(jsonConfig)
		if err != nil {
			l.Logger.Error().Msg(err.Error())
		}

		cacheConfig.Log = &l.Logger
		cacheStorage, err := cache.MakeStorage(cacheConfig)
		if err != nil {
			l.Logger.Error().Msg(err.Error())
		}

		/* configure and start API server */
		apiConfig.Log = &l.Logger
		err = apiConfig.ConfigureFromInterface(jsonConfig)

		if err != nil {
			l.Logger.Error().Msg(err.Error())
		}

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
