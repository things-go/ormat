package tool

import (
	"errors"
	stdlog "log"
	"os"
	"os/exec"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/things-go/ormat/database"
	"github.com/things-go/ormat/deploy"
	"github.com/things-go/ormat/log"
	"github.com/things-go/ormat/utils"
	"github.com/things-go/ormat/view"
	"github.com/things-go/ormat/view/driver"
)

// Execute exe the cmd
func Execute() {
	db, md, err := GetDbAndViewModel()
	if err != nil {
		log.Error(err)
		return
	}
	defer database.Close(db)

	cfg := GetConfig()

	_, dbName, err := cfg.GetDbDSNAndDbName()
	if err != nil {
		log.Error(err)
		return
	}

	vw := view.New(db, md, cfg.View, dbName, cfg.TableNames...)

	list, err := vw.GetDbFile(utils.GetPkgName(cfg.OutDir))
	if err != nil {
		log.Error(err)
		return
	}
	for _, v := range list {
		path := cfg.OutDir + "/" + v.GetName()
		_ = utils.WriteFile(path, []byte(v.Build()))

		cmd, _ := exec.Command("goimports", "-l", "-w", path).Output()
		log.Info(strings.TrimSuffix(string(cmd), "\n"))

		_, _ = exec.Command("gofmt", "-l", "-w", path).Output()
	}

	log.Info("generate success !!!")
}

func ExecuteCreateSQL() {
	db, md, err := GetDbAndViewModel()
	if err != nil {
		log.Error(err)
		return
	}
	defer database.Close(db)
	cfg := GetConfig()
	_, dbName, err := cfg.GetDbDSNAndDbName()
	if err != nil {
		log.Error(err)
		return
	}

	vw := view.New(db, md, cfg.View, dbName, cfg.TableNames...)

	content, err := vw.GetDBCreateTableSQLContent()
	if err != nil {
		log.Error(err)
		return
	}
	_ = utils.WriteFile(cfg.OutDir+"/create_table.sql", content)

}

func GetDbAndViewModel() (*gorm.DB, view.DBModel, error) {
	var gc = &gorm.Config{}

	if !deploy.IsRelease() {
		gc.Logger = logger.New(stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		})
	}
	cfg := GetConfig()

	dsn, _, err := cfg.GetDbDSNAndDbName()
	if err != nil {
		return nil, nil, err
	}
	switch cfg.GetDatabase().Dialect {
	case "mysql": // mysql
		db, err := database.New(database.Config{Dialect: "mysql", Dsn: dsn}, gc)
		return db, &driver.MySQL{CustomDefineType: cfg.GetTypeDefine()}, err
	case "sqlite3": // sqlite3
		db, err := database.New(database.Config{Dialect: "sqlite3", Dsn: dsn}, gc)
		return db, &driver.SQLite{CustomDefineType: cfg.GetTypeDefine()}, err
	default:
		return nil, nil, errors.New("database not fund: please check database.dialect (mysql, sqlite3, mssql)")
	}
}
