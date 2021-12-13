package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thinkgos/ormat/pkg/infra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Get version info",
	Example: "ormat version",
	RunE: func(*cobra.Command, []string) error {
		fmt.Println("ormat " + infra.Version)
		return nil
	},
}
