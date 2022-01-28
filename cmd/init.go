package cmd

import (
	"github.com/spf13/cobra"

	"github.com/thinkgos/ormat/pkg/infra"
	"github.com/thinkgos/ormat/tpl"
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
		return infra.WriteFile(".ormat.yml", b)
	},
}
