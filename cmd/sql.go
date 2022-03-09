package cmd

import (
	"github.com/spf13/cobra"

	"github.com/things-go/ormat/tool"
)

var sqlCmd = &cobra.Command{
	Use:     "sql",
	Short:   "generate create table sql",
	Example: "ormat sql",
	RunE: func(*cobra.Command, []string) error {
		initConfig()
		tool.ExecuteCreateSQL()
		return nil
	},
}
