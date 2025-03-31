package cmd

import (
	"github.com/apsamuel/brainiac/pkg/logger"
	"github.com/spf13/cobra"
)

var apiCommand = &cobra.Command{
	Use:   "api",
	Short: "api commands",
	Long:  "api commands",
	Run: func(cmd *cobra.Command, args []string) {
		l := logger.Logger
		l.Logger.Info().Msg("AI command called")
	},
}
