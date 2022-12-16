package cmd

import (
	"github.com/spf13/cobra"
	"github.com/things-go/log"

	"github.com/things-go/ormat/pkg/consts"
)

type RootCmd struct {
	cmd        *cobra.Command
	configFile string
	level      string
}

func NewRootCmd() *RootCmd {
	root := &RootCmd{}
	cmd := &cobra.Command{
		Use:           "ormat",
		Short:         "gorm reflect tools",
		Long:          "database/sql to golang struct",
		Version:       consts.BuildVersion(),
		SilenceUsage:  false,
		SilenceErrors: false,
		Args:          cobra.NoArgs,
	}
	cobra.OnInitialize(func() {
		log.ReplaceGlobals(log.NewLogger(log.WithConfig(log.Config{
			Level:  root.level,
			Format: "console",
		})))
	})

	cmd.PersistentFlags().StringVarP(&root.configFile, "config", "c", "", "config file")
	cmd.PersistentFlags().StringVarP(&root.level, "level", "l", "info", "log level(debug,info,warn,error,dpanic,panic,fatal)")
	cmd.AddCommand(
		newSqlCmd().cmd,
		newBuildCmd().cmd,
		newGenCmd().cmd,
		newExpandCmd().cmd,
		newUpgradeCmd().cmd,
	)
	root.cmd = cmd
	return root
}

// Execute adds all child commands to the root command and sets flags appropriately.
func (r *RootCmd) Execute() error {
	return r.cmd.Execute()
}
