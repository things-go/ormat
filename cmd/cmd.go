package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var outDir string

var rootCmd = &cobra.Command{
	Use:   "ormat",
	Short: "gorm reflect tools",
	Long:  "database/sql to golang struct",
}

func init() {
	rootCmd.AddCommand(
		versionCmd,
		configCmd,
		sqlCmd,
		buildCmd,
		genCmd,
		expandCmd,
	)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
