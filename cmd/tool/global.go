package tool

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/things-go/ormat/config"
	"github.com/things-go/ormat/deploy"
	"github.com/things-go/ormat/utils"
	"github.com/things-go/ormat/view"
)

var cfg = config.Config{
	Deploy: deploy.Prod.String(),
	Database: config.Database{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "root",
		Password: "root",
		Db:       "test",
	},
	OutDir:     "./model",
	TypeDefine: make(map[string]string),
	View: view.Config{
		DbTag:           "gorm",
		WebTags:         []view.WebTag{{Kind: view.WebTagSnakeCase, Tag: "json", HasOmit: true}},
		EnableLint:      false,
		DisableNull:     false,
		EnableInt:       false,
		IsNullToPoint:   true,
		IsOutSQL:        false,
		IsOutColumnName: false,
		IsForeignKey:    false,
		IsCommentTag:    false,
	},
}

func LoadConfig() error {
	viper.SetConfigName(".ormat")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(utils.ExecutableDir())
	viper.AddConfigPath(utils.WorkDir())
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&cfg, func(c *mapstructure.DecoderConfig) { c.TagName = "yaml" })
	if err != nil {
		return err
	}
	return cfg.Database.Parse()
}

// GetConfig get config
func GetConfig() config.Config {
	return cfg
}
