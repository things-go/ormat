package cmd

import (
	"github.com/spf13/cobra"

	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/pkg/tpl"
	"github.com/things-go/ormat/pkg/utils"
)

func init() {

}

type configCmd struct {
	cmd *cobra.Command
}

func newConfigCmd() *configCmd {
	root := &configCmd{}
	cmd := &cobra.Command{
		Use:     "config",
		Short:   "Show/Generate config file",
		Example: "ormat config - show config \normat config init - Generate config file",
		RunE: func(*cobra.Command, []string) error {
			c := config.Global
			err := c.Load()
			if err != nil {
				return err
			}
			JSON(c)
			return nil
		},
	}

	cmdInit := &cobra.Command{
		Use:     "init",
		Short:   "Generate config file",
		Example: "ormat config init",
		RunE: func(*cobra.Command, []string) error {
			b, err := tpl.Static.ReadFile("template/ormat.yml")
			if err != nil {
				return err
			}
			return utils.WriteFile(".ormat.yml", b)
		},
	}
	cmd.AddCommand(cmdInit)
	root.cmd = cmd
	return root
}
