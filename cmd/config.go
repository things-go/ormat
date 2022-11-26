package cmd

import (
	"encoding/json"
	"os"

	"github.com/alecthomas/chroma/quick"
	"github.com/spf13/cobra"

	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/pkg/tpl"
	"github.com/things-go/ormat/pkg/utils"
)

func init() {
	configCmd.AddCommand(configInitSubCmd)
}

var configCmd = &cobra.Command{
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

var configInitSubCmd = &cobra.Command{
	Use:     "init",
	Short:   "Generate config file",
	Example: "ormat config init",
	RunE: func(*cobra.Command, []string) error {
		b, err := tpl.Static.ReadFile("ormat.yml")
		if err != nil {
			return err
		}
		return utils.WriteFile(".ormat.yml", b)
	},
}

func JSON(v ...interface{}) {
	for _, vv := range v {
		b, _ := json.MarshalIndent(vv, "", "  ")
		quick.Highlight(os.Stdout, string(b), "JSON", "terminal", "solarized-dark") // nolint
	}
}
