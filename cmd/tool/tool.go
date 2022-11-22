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

	"github.com/things-go/log"

	"github.com/things-go/ormat/database"
	"github.com/things-go/ormat/deploy"
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

	c := GetConfig()
	vw := view.New(md, c.View)

	list, err := vw.GetDbFile(utils.GetPkgName(c.OutDir))
	if err != nil {
		log.Error(err)
		return
	}
	for _, v := range list {
		path := c.OutDir + "/" + v.Filename + ".go"
		_ = utils.WriteFile(path, []byte(v.Build()))

		cmd, _ := exec.Command("goimports", "-l", "-w", path).Output()
		log.Info(strings.TrimSuffix(string(cmd), "\n"))
		_, _ = exec.Command("gofmt", "-l", "-w", path).Output()

		if c.View.IsOutSQL {
			_ = utils.WriteFile(c.OutDir+"/"+v.Filename+".sql", []byte(v.BuildSQL()))
		}
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
	c := GetConfig()

	vw := view.New(md, c.View)

	content, err := vw.GetDBCreateTableSQLContent()
	if err != nil {
		log.Error(err)
		return
	}
	_ = utils.WriteFile(c.OutDir+"/create_table.sql", content)
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
	c := GetConfig()
	dbCfg := c.Database

	switch dbCfg.Dialect {
	case "mysql": // mysql
		db, err := database.New(database.Config{Dialect: "mysql", Dsn: dbCfg.Dsn()}, gc)
		if err != nil {
			return nil, nil, err
		}
		return db, &driver.MySQL{
			DB:               db,
			DbName:           dbCfg.DbName(),
			TableNames:       c.TableNames,
			CustomDefineType: c.TypeDefine,
		}, err
	case "sqlite3": // sqlite3
		db, err := database.New(database.Config{Dialect: "sqlite3", Dsn: dbCfg.Dsn()}, gc)
		if err != nil {
			return nil, nil, err
		}
		return db, &driver.SQLite{
			DB:               db,
			DbName:           dbCfg.DbName(),
			TableNames:       c.TableNames,
			CustomDefineType: c.TypeDefine,
		}, nil
	default:
		return nil, nil, errors.New("database not fund: please check database.dialect (mysql, sqlite3, mssql)")
	}
}
