package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/things-go/ormat/pkg/utils"
)

var rootCmd = &cobra.Command{
	Use:   "ormat",
	Short: "gorm reflect tools",
	Long:  "database/sql to golang struct",
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(
		versionCmd,
		configCmd,
		sqlCmd,
		buildCmd,
		genCmd,
		expandCmd,
	)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetConfigName(".ormat")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(utils.WorkDir())

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
