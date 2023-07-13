package command

import (
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	mysqlib "github.com/go-sql-driver/mysql"
	"golang.org/x/exp/slog"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DbConfig connect information
type DbConfig struct {
	Dialect string `yaml:"dialect" json:"dialect" binding:"oneof=mysql sqlite3"` // mysql, sqlite3
	DSN     string `yaml:"dsn" json:"dsn"`
	Options string `yaml:"options" json:"options"` // Options ?号后面, 如果为空, 则为 charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=True
}

func NewDB(c DbConfig) (*gorm.DB, string, error) {
	var dialector gorm.Dialector
	var dbName string

	switch c.Dialect {
	case "mysql":
		dsn := c.DSN
		if c.Options == "" {
			c.Options = "charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=True"
		}
		idx := strings.Index(dsn, "?")
		if idx == -1 {
			dsn = fmt.Sprintf("%s?%s", dsn, c.Options)
		}
		cc, err := mysqlib.ParseDSN(dsn)
		if err != nil {
			return nil, "", err
		}
		dbName = cc.DBName
		dialector = mysql.New(mysql.Config{
			DSN: dsn,
			// DefaultStringSize:         256,   // string 类型字段的默认长度
			// DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			// DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			// DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			// SkipInitializeWithVersion: false, // 根据版本自动配置
		})
	case "postgres":
		dialector = postgres.New(postgres.Config{DSN: c.DSN})
	case "sqlite3":
		dsn := c.DSN

		dbName = filepath.Base(dsn)
		if dbName == "" {
			return nil, "", errors.New("empty sqlite3 db name")
		}

		// 路径是否存在
		_, err := os.Stat(dsn)
		if !(err == nil || os.IsExist(err)) {
			if err := os.MkdirAll(path.Dir(dsn), os.ModePerm); err != nil {
				return nil, "", fmt.Errorf("database mkdir (%s), %+v", dsn, err)
			}
			if _, err := os.Create(dsn); err != nil {
				return nil, "", fmt.Errorf("database create DB(%s), %+v", dsn, err)
			}
		}
		dialector = NewSqlite3(dsn)
	default:
		return nil, "", errors.New("please select database driver one of [mysql|postgres|sqlite3], if use sqlite3, build tags with sqlite3.")
	}
	config := &gorm.Config{}

	if slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		config.Logger = logger.New(stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), logger.Config{
			SlowThreshold:             500 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		})
	}
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, "", err
	}
	return db, dbName, nil
}

func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
