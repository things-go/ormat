package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/things-go/ormat/pkg/consts"
	"github.com/things-go/ormat/pkg/utils"
)

type RootCmd struct {
	cmd        *cobra.Command
	configFile string
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
		if root.configFile != "" {
			viper.SetConfigFile(root.configFile)
		} else {
			viper.AddConfigPath(utils.WorkDir())
			viper.SetConfigName(".ormat")
			viper.SetConfigType("yaml")
		}
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	})

	cmd.PersistentFlags().StringVarP(&root.configFile, "config", "c", "", "config file")
	cmd.AddCommand(
		newConfigCmd().cmd,
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
