package cmd

import (
	"os"

	"github.com/apsamuel/brainiac/pkg/common"
	"github.com/apsamuel/brainiac/pkg/logger"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

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
