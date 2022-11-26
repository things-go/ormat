package cmd

import (
	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/runtime"
	"github.com/things-go/ormat/view"
	"github.com/things-go/ormat/view/driver"
)

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
