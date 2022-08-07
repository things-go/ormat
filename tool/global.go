package tool

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/thinkgos/ormat/config"
	"github.com/thinkgos/ormat/deploy"
	"github.com/thinkgos/ormat/utils"
	"github.com/thinkgos/ormat/view"
)

var cfg = config.Config{
	Deploy: deploy.Prod,
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
		DbTag:         "gorm",
		WebTags:       []view.WebTag{{Kind: "snakeCase", Tag: "json", HasOmit: true}},
		EnableLint:    false,
		DisableNull:   false,
		EnableInt:     false,
		IsNullToPoint: true,
		IsOutSQL:      false,
		IsForeignKey:  false,
		IsCommentTag:  false,
	},
}

func LoadConfig() error {
	viper.SetConfigName(".ormat")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(utils.GetExecutableDir())
	viper.AddConfigPath(utils.GetWd())
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return viper.Unmarshal(&cfg, func(c *mapstructure.DecoderConfig) { c.TagName = "yaml" })
}

// GetConfig get config
func GetConfig() config.Config { return cfg }
