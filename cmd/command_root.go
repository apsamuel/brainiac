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
	// TODO: understand why configPort is never set
	rootCommand.PersistentFlags().IntVarP(&configPort, "configPort", "p", 5432, "configuration port")
	rootCommand.PersistentFlags().StringVarP(&configDatabase, "configDatabase", "D", "observability", "configuration database")
	rootCommand.PersistentFlags().StringVarP(&configTable, "configTable", "t", "config_data", "configuration table or key")
	rootCommand.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debugging")

	rootCommand.AddCommand(runCommand)
	// rootCommand.AddCommand(runxCommand)
	rootCommand.AddCommand(configCommand)
	rootCommand.AddCommand(aiCommand)
	rootCommand.AddCommand(apiCommand)
}
