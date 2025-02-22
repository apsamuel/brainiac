package cmd

import (
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "brainiac",
	Short: "Brainiac is a tool for AI",
}

func Execute() error {
	return rootCommand.Execute()
}

func init() {
	rootCommand.PersistentFlags().StringVarP(&configFile, "configYaml", "c", "", "config file")
	rootCommand.PersistentFlags().StringVarP(&configHost, "configHost", "H", "localhost", "host")
	rootCommand.PersistentFlags().IntVarP(&configPort, "configPort", "p", 8080, "port")
	rootCommand.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug")

	rootCommand.AddCommand(runCommand)
	rootCommand.AddCommand(configCommand)
}
