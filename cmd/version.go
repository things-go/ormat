package cmd

import (
	"github.com/spf13/cobra"

	"github.com/thinkgos/ormat/pkg/builder"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Get version info",
	Example: "ormat version",
	RunE: func(*cobra.Command, []string) error {
		builder.PrintVersion()
		return nil
	},
}
