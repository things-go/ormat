package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "v0.0.1-rc5"

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Get version info",
	Example: "ormat version",
	RunE: func(*cobra.Command, []string) error {
		fmt.Println("ormat " + version)
		return nil
	},
}
