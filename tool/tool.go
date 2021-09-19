package tool

import (
	"errors"
	stdlog "log"
	"os"
	"os/exec"
	"time"

	"github.com/things-go/x/extos"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/thinkgos/ormat/config"
	"github.com/thinkgos/ormat/pkg/database"
	"github.com/thinkgos/ormat/pkg/env"
	"github.com/thinkgos/ormat/pkg/infra"
	"github.com/thinkgos/ormat/pkg/zapl"
	"github.com/thinkgos/ormat/view"
	"github.com/thinkgos/ormat/view/driver"
)

// Execute exe the cmd
func Execute() {
	db, md, err := GetDbAndViewModel()
	if err != nil {
		zapl.Error(err)
		return
	}
	defer database.Close(db)

	_, dbName := config.GetDatabaseDSNAndDbName()

	cfg := config.GetConfig()
	vw := view.New(db, md, cfg.View, dbName, cfg.TableNames...)

	list, err := vw.GetDbFile(infra.GetPkgName(cfg.OutDir))
	if err != nil {
		zapl.Error(err)
		return
	}
	for _, v := range list {
		path := cfg.OutDir + "/" + v.GetName()
		_ = extos.WriteFile(path, []byte(v.Build()))

		zapl.Info("run goimports")
		cmd, _ := exec.Command("goimports", "-l", "-w", path).Output()
		zapl.Info(string(cmd))

		zapl.Info("run gofmt")
		cmd, _ = exec.Command("gofmt", "-l", "-w", path).Output()
		zapl.Info(string(cmd))
	}

	zapl.Info("generate success !!!")
}

func ExecuteCreateSQL() {
	db, md, err := GetDbAndViewModel()
	if err != nil {
		zapl.Error(err)
		return
	}
	defer database.Close(db)

	_, dbName := config.GetDatabaseDSNAndDbName()
	cfg := config.GetConfig()
	vw := view.New(db, md, cfg.View, dbName, cfg.TableNames...)

	content, err := vw.GetDBCreateTableSQLContent()
	if err != nil {
		zapl.Error(err)
		return
	}
	_ = extos.WriteFile(cfg.OutDir+"/create_table.sql", content)

}

func GetDbAndViewModel() (*gorm.DB, view.DBModel, error) {
	var gc = &gorm.Config{}

	if !env.IsDeployRelease() {
		gc.Logger = logger.New(stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		})
	}

	dsn, _ := config.GetDatabaseDSNAndDbName()
	switch config.GetDbInfo().Dialect {
	case "mysql": // mysql
		db, err := database.New(database.Config{Dialect: "mysql", Dsn: dsn}, gc)
		return db, &driver.MySQL{}, err
	case "sqlite3": // sqlite3
		db, err := database.New(database.Config{Dialect: "sqlite3", Dsn: dsn}, gc)
		return db, &driver.SQLite{}, err
	default:
		return nil, nil, errors.New("database not fund: please check database.dialect (mysql, sqlite3, mssql)")
	}
}
