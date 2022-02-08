package cmd

import (
	"fmt"
	"os"

	validator "github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"

	"github.com/thinkgos/ormat/config"
	"github.com/thinkgos/ormat/pkg/env"
	"github.com/thinkgos/ormat/pkg/zapl"
	"github.com/thinkgos/ormat/tool"
)

var validate = validator.New()

var rootCmd = &cobra.Command{
	Use:   "ormat",
	Short: "gorm reflect tools",
	Long:  "database to golang struct",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		tool.Execute()
	},
}

func init() {
	validate.SetTagName("binding")

	rootCmd.AddCommand(versionCmd, initCmd, sqlCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// initConfig reads in config file.
func initConfig() {
	zapl.ReplaceGlobals(zapl.New(zapl.Config{Level: "info", Format: "console"}).Sugar())
	err := config.LoadConfig()
	if err != nil {
		zapl.Fatalf("load config failed(please run 'ormat init' generate a .ormat.yml): %s", err.Error())
		return
	}
	c := config.GetConfig()
	err = validate.Struct(c)
	if err != nil {
		zapl.Info("config validate failed: use（-h, --help) to get more info")
		zapl.Error(err)
		os.Exit(1)
		return
	}
	env.SetDeploy(c.Deploy)
	fmt.Println("using config:")
	zapl.JSON(c)
}
