package cmd

import (
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "brainiac",
	Short: "brainiac is a tool for AI",
}

func Execute() error {
	return rootCommand.Execute()
}

func init() {
	rootCommand.PersistentFlags().StringVarP(&configEngine, "configEngine", "e", "postgres", "configuration engine")
	rootCommand.PersistentFlags().StringVarP(&configFile, "configYaml", "c", "", "configuration file")
	// TODO: use host:port string instead of separate flags
	rootCommand.PersistentFlags().StringVarP(&configHost, "configHost", "H", "localhost", "configuration host")
	rootCommand.PersistentFlags().IntVarP(&configPort, "configPort", "p", 8080, "configuration port")
	rootCommand.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debugging")

	rootCommand.AddCommand(runCommand)
	rootCommand.AddCommand(configCommand)
}
