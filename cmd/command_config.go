package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/apsamuel/brainiac/pkg/cache"
	"github.com/apsamuel/brainiac/pkg/common"
	"github.com/apsamuel/brainiac/pkg/database"
	"github.com/apsamuel/brainiac/pkg/logger"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

var PostgresClient *gorm.DB

var supportedEngines = []string{"postgres", "redis"}

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "manage brainiac configuration",
	Long:  `manage brainiac configuration`,
	Run: func(cmd *cobra.Command, args []string) {

		l := logger.Logger
		configInterface := make(map[string]interface{})
		data, err := os.ReadFile(configFile)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(data, &configInterface)
		if err != nil {
			panic(err)
		}

		if generateSecret {
			key, nonce, err := common.GenerateSecret(keySize, nonceSize)
			if err != nil {
				panic(err)
			}
			l.Logger.Info().Str("key", key).Str("nonce", nonce).Msg("generated secret")
			return
		}

		if configWrite {
			if !common.Contains(supportedEngines, configEngine) {
				panic(fmt.Sprintf("unsupported configuration engine: %s", configEngine))
			}

			l.Logger.Info().Msgf("writing configuration to %s", configEngine)
			if configEngine == "postgres" {

				configHost := os.Getenv("BRAINIAC_CONFIG_HOST")
				if configHost == "" {
					panic("BRAINIAC_CONFIG_HOST not set")
				}
				configPortStr := os.Getenv("BRAINIAC_CONFIG_PORT")
				if configPortStr == "" {
					panic("BRAINIAC_CONFIG_PORT not set")
				} else {
					configPort, err = strconv.Atoi(configPortStr)
					if err != nil {
						panic(err)
					}
				}
				configDatabase := os.Getenv("BRAINIAC_CONFIG_DB")
				if configDatabase == "" {
					panic("BRAINIAC_CONFIG_DB not set")
				}
				configUsername := os.Getenv("BRAINIAC_CONFIG_USER")
				if configUsername == "" {
					panic("BRAINIAC_CONFIG_USER not set")
				}
				configPassword := os.Getenv("BRAINIAC_CONFIG_PASS")
				if configPassword == "" {
					panic("BRAINIAC_CONFIG_PASS not set")
				}
				configKey := os.Getenv(aesKeyVariable)
				if configKey == "" {
					panic(fmt.Sprintf("%s not set", aesKeyVariable))
				}
				configNonce := os.Getenv(aesNonceVariable)
				if configNonce == "" {
					panic(fmt.Sprintf("%s not set", aesNonceVariable))
				}

				database.PushConfig(
					configHost,
					configPort,
					configDatabase,
					"config_data",
					configUsername,
					configPassword,
					data,
					configKey,
					configNonce,
				)
			}

			if configEngine == "redis" {
				configHost := os.Getenv("BRAINIAC_CONFIG_HOST")
				if configHost == "" {
					panic("BRAINIAC_CONFIG_HOST not set")
				}
				configPortStr := os.Getenv("BRAINIAC_CONFIG_PORT")
				if configPortStr == "" {
					panic("BRAINIAC_CONFIG_PORT not set")
				} else {
					configPort, err = strconv.Atoi(configPortStr)
					if err != nil {
						panic(err)
					}
				}
				configKey := os.Getenv(aesKeyVariable)
				if configKey == "" {
					panic("BRAINIAC_AES_KEY not set")
				}
				configNonce := os.Getenv(aesNonceVariable)
				if configNonce == "" {
					panic("BRAINIAC_AES_NONCE not set")
				}
				cipherText, err := common.EncryptWithAESGCM(data, []byte(configKey))
				if err != nil {
					panic(err)
				}

				cache.PushConfig(
					configHost,
					configPort,
					"config_data",
					[]byte(cipherText),
				)
			}

			return
		}
	},
}

func init() {
	configCommand.Flags().BoolVarP(&configWrite, "writeConfig", "w", false, "write configuration")
	configCommand.Flags().BoolVarP(&generateSecret, "getSecret", "g", false, "generate brainiac secret")
}
