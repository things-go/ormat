package cmd

import (
	"github.com/spf13/cobra"

	"github.com/thinkgos/ormat/cmd/tpl"
	"github.com/thinkgos/ormat/utils"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "generate config file",
	Example: "ormat init",
	RunE: func(*cobra.Command, []string) error {
		b, err := tpl.Static.ReadFile("ormat.yml")
		if err != nil {
			return err
		}
		return utils.WriteFile(".ormat.yml", b)
	},
}
