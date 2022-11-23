package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ormat",
	Short: "gorm reflect tools",
	Long:  "database to golang struct",
}

func init() {
	rootCmd.AddCommand(
		versionCmd,
		initCmd,
		sqlCmd,
		dbCmd,
		genCmd,
	)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
