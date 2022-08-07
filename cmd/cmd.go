package cmd

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"

	"github.com/things-go/ormat/deploy"
	"github.com/things-go/ormat/log"
	"github.com/things-go/ormat/tool"
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
	log.ReplaceGlobals(log.New(log.Config{Level: "info", Format: "console"}).Sugar())
	err := tool.LoadConfig()
	if err != nil {
		log.Fatalf("load config failed(please run 'ormat init' generate a .ormat.yml): %s", err.Error())
		return
	}
	c := tool.GetConfig()
	err = validate.Struct(c)
	if err != nil {
		log.Info("config validate failed: use（-h, --help) to get more info")
		log.Error(err)
		os.Exit(1)
		return
	}
	deploy.Set(c.Deploy)
	fmt.Println("using config:")
	log.JSON(c)
}
