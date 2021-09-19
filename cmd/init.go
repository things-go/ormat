package cmd

import (
	"github.com/spf13/cobra"
	"github.com/things-go/x/extos"

	"github.com/thinkgos/ormat/tpl"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "generate config file",
	Example: "ormat init",
	RunE: func(*cobra.Command, []string) error {
		b, err := tpl.Static.ReadFile("config.yml")
		if err != nil {
			return err
		}
		return extos.WriteFile("config.yml", b)
	},
}
