package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/things-go/ormat/consts"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Get version info",
	Example: "ormat version",
	RunE: func(*cobra.Command, []string) error {
		fmt.Println("ormat " + consts.Version)
		return nil
	},
}
