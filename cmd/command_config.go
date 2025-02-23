package cmd

import (
	"os"

	"github.com/apsamuel/brainiac/pkg/common"
	"github.com/apsamuel/brainiac/pkg/logger"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var supportedEngines = []string{"postgres", "redis"}

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "Brainiac configuration",
	Long:  `Brainiac configuration`,
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

		if configDebug {
			l.Logger.Debug().Msgf("debugging configuration %v", configInterface)
			return
		}

		if configWrite {
			configEngine := os.Getenv("BRAINIAC_CONFIG_ENGINE")
			if configEngine == "" {
				panic("BRAINIAC_CONFIG_ENGINE not set")
			}
			if !common.Contains(supportedEngines, configEngine) {
				panic("unsupported engine")
			}
			configHost := os.Getenv("BRAINIAC_CONFIG_HOST")
			if configHost == "" {
				panic("BRAINIAC_CONFIG_HOST not set")
			}
			configPort := os.Getenv("BRAINIAC_CONFIG_PORT")
			if configPort == "" {
				panic("BRAINIAC_CONFIG_PORT not set")
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
			configKey := os.Getenv("BRAINIAC_AES_KEY")
			if configKey == "" {
				panic("BRAINIAC_AES_KEY not set")
			}
			configNonce := os.Getenv("BRAINIAC_AES_NONCE")
			if configNonce == "" {
				panic("BRAINIAC_AES_NONCE not set")
			}
			cipherText, err := common.EncryptWithAESGCM(data, []byte(configKey))
			if err != nil {
				panic(err)
			}

			l.Logger.Info().Str("ciperText", cipherText).Msg("debugging encrypted configuration")
			return
		}
	},
}

func init() {
	configCommand.Flags().BoolVarP(&configDebug, "debug", "d", false, "debug configuration")
	configCommand.Flags().BoolVarP(&configWrite, "write", "w", false, "write configuration")
	configCommand.Flags().BoolVarP(&generateSecret, "generate-secret", "g", false, "generate brainiac secret")
}
