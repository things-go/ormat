package cmd

import (
	"github.com/go-playground/validator/v10"
	"github.com/things-go/log"
	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/pkg/deploy"
	"github.com/things-go/ormat/runtime"
	"github.com/things-go/ormat/view"
	"github.com/things-go/ormat/view/driver"
)

var validate = validator.New()

func init() {
	validate.SetTagName("binding")
}

func GetViewModel(rt *runtime.Runtime) view.DBModel {
	c := config.Global
	dbCfg := c.Database
	switch dbCfg.Dialect {
	case "mysql":
		return &driver.MySQL{
			DB:               rt.DB,
			DbName:           dbCfg.DbName(),
			TableNames:       c.TableNames,
			CustomDefineType: c.TypeDefine,
		}
	case "sqlite3":
		return &driver.SQLite{
			DB:               rt.DB,
			DbName:           dbCfg.DbName(),
			TableNames:       c.TableNames,
			CustomDefineType: c.TypeDefine,
		}
	default:
		panic("database not found, please check database.dialect (mysql, sqlite3, mssql)")
	}
}

func setupBase(c *config.Config) {
	deploy.MustSetDeploy(c.Deploy)
	log.ReplaceGlobals(log.NewLogger(log.WithConfig(log.Config{
		Level:  "info",
		Format: "console",
	})))
}
