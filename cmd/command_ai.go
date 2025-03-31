package cmd

import (
	"github.com/apsamuel/brainiac/pkg/logger"
	"github.com/spf13/cobra"
)

var aiCommand = &cobra.Command{
	Use:   "ai",
	Short: "ai commands",
	Long:  "ai commands",
	Run: func(cmd *cobra.Command, args []string) {
		l := logger.Logger
		l.Logger.Info().Msg("AI command called")
	},
}
