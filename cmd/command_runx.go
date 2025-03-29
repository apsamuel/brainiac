package cmd

import "github.com/spf13/cobra"

/*
Start each extension as it's own indepdendent service
*/

var runxCommand = &cobra.Command{
	Use:   "runx",
	Short: "Run the extension",
	Long:  "Start the extension as an independent service.",
}

func init() {
	runxCommand.Flags().StringVarP(&configFile, "configFile", "f", "", "configuration file")
	runxCommand.Flags().StringVarP(&configEngine, "configEngine", "e", "", "configuration engine")
	runxCommand.Flags().StringVarP(&configHost, "configHost", "H", "", "configuration host")
	runxCommand.Flags().IntVarP(&configPort, "configPort", "p", 0, "configuration port")
	runxCommand.Flags().StringVarP(&configDatabase, "configDatabase", "D", "", "configuration database")
	runxCommand.Flags().StringVarP(&configTable, "configTable", "t", "", "configuration table or key")
	rootCommand.AddCommand(runxCommand)
}
