package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/things-go/ormat/pkg/consts"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Show version",
	Example: "ormat version",
	RunE: func(*cobra.Command, []string) error {
		fmt.Println("ormat version " + consts.Version)
		return nil
	},
}
