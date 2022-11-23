package runtime

import (
	"encoding/json"
	"fmt"
	stdlog "log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/things-go/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/pkg/database"
	"github.com/things-go/ormat/pkg/deploy"
	"github.com/things-go/ormat/pkg/utils"
)

var ConfigSet = wire.NewSet(NewConfig)

func NewConfig(remote bool) (*config.Config, error) {
	c := config.NewDefaultConfig()

	doErr := func(e error) (*config.Config, error) {
		if remote {
			return nil, e
		}
		// 如果不是远程, 则使用默认
		return &c, nil
	}

	viper.SetConfigName(".ormat")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(utils.ExecutableDir())
	viper.AddConfigPath(utils.WorkDir())
	err := viper.ReadInConfig()
	if err != nil {
		return doErr(err)
	}
	err = viper.Unmarshal(&c, func(c *mapstructure.DecoderConfig) { c.TagName = "yaml" })
	if err != nil {
		return doErr(err)
	}
	if remote {
		err = c.Database.Parse()
		if err != nil {
			return nil, err
		}
	}
	validate := validator.New()
	validate.SetTagName("binding")
	err = validate.Struct(c)
	if err != nil {
		return nil, err
	}
	deploy.MustSetDeploy(c.Deploy)
	log.ReplaceGlobals(log.NewLogger(log.WithConfig(log.Config{Level: "info", Format: "console"})))
	JSON(c)
	return &c, nil
}

var DbSet = wire.NewSet(NewDb, NewDbConfig)

func NewDb(remote bool, dbCfg config.Database) (*gorm.DB, error) {
	if remote {
		gc := &gorm.Config{}
		if !deploy.IsRelease() {
			gc.Logger = logger.New(stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
			})
		}

		return database.New(database.Config{
			Dialect: dbCfg.Dialect,
			Dsn:     dbCfg.Dsn(),
		}, gc)
	}
	return nil, nil
}

func NewDbConfig(c *config.Config) config.Database {
	return c.Database
}

func JSON(v ...interface{}) {
	for _, vv := range v {
		b, _ := json.MarshalIndent(vv, "", "  ")
		fmt.Println(string(b))
	}
}
