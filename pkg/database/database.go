package database

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// Node 数据库配置
type Node struct {
	// dsn
	Dsn string `yaml:"dsn" json:"dsn"`
}

type SourceServe struct {
	// MaxIdleConn sets the maximum number of open connections to the database.
	// If n <= 0, then there is no limit on the number of open connections.
	// The default is 0 (unlimited).
	MaxIdleConn int `yaml:"maxIdleConn" json:"maxIdleConn"`
	// MaxOpenConn sets the maximum number of open connections to the database.
	// If n <= 0, then there is no limit on the number of open connections.
	// The default is 0 (unlimited).
	MaxOpenConn int `yaml:"maxOpenConn" json:"maxOpenConn"`
	// MaxLifetime sets the maximum amount of time a connection may be reused.
	// If d <= 0, connections are not closed due to a connection's age.
	MaxLifetime time.Duration `yaml:"maxLifetime" json:"maxLifetime"`
	// MaxIdleTime sets the maximum amount of time a connection may be idle.
	// If d <= 0, connections are not closed due to a connection's idle time.
	MaxIdleTime time.Duration `yaml:"maxIdleTime" json:"maxIdleTime"`
	Node        Node          `yaml:"node" json:"node"`
}

type ReplicaServe struct {
	// MaxIdleConn sets the maximum number of open connections to the database.
	// If n <= 0, then there is no limit on the number of open connections.
	// The default is 0 (unlimited).
	MaxIdleConn int `yaml:"maxIdleConn" json:"maxIdleConn"`
	// MaxOpenConn sets the maximum number of open connections to the database.
	// If n <= 0, then there is no limit on the number of open connections.
	// The default is 0 (unlimited).
	MaxOpenConn int `yaml:"maxOpenConn" json:"maxOpenConn"`
	// MaxLifetime sets the maximum amount of time a connection may be reused.
	// If d <= 0, connections are not closed due to a connection's age.
	MaxLifetime time.Duration `yaml:"maxLifetime" json:"maxLifetime"`
	// MaxIdleTime sets the maximum amount of time a connection may be idle.
	// If d <= 0, connections are not closed due to a connection's idle time.
	MaxIdleTime time.Duration `yaml:"maxIdleTime" json:"maxIdleTime"`
	Node        []Node        `yaml:"node" json:"node"`
}

type Config struct {
	// EnableLog enabled log flag  use by user
	EnableLog bool `yaml:"enableLog" json:"enableLog"`
	// mysql sqlite3 postgres
	Dialect string       `yaml:"dialect" json:"dialect"`
	Source  SourceServe  `yaml:"source" json:"source"`
	Replica ReplicaServe `yaml:"replica" json:"replica"`
}

func New(c *Config, config *gorm.Config) (*gorm.DB, error) {
	sourceDialector, err := newDialector(c.Dialect, &c.Source.Node)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(sourceDialector, config)
	if err != nil {
		return nil, err
	}
	// replica
	if len(c.Replica.Node) > 0 {
		replicaDialector := make([]gorm.Dialector, 0, len(c.Replica.Node))
		for _, node := range c.Replica.Node {
			dialect, err := newDialector(c.Dialect, &node)
			if err != nil {
				return nil, err
			}
			replicaDialector = append(replicaDialector, dialect)
		}
		pluginDbResolver := dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{sourceDialector},
			Replicas: replicaDialector,
			Policy:   nil,
		})
		if c.Source.MaxIdleConn > 0 {
			pluginDbResolver.SetMaxIdleConns(c.Replica.MaxIdleConn)
		}
		if c.Source.MaxOpenConn > 0 {
			pluginDbResolver.SetMaxOpenConns(c.Replica.MaxOpenConn)
		}
		if c.Source.MaxLifetime > 0 {
			pluginDbResolver.SetConnMaxLifetime(c.Replica.MaxLifetime)
		}
		if c.Source.MaxIdleTime > 0 {
			pluginDbResolver.SetConnMaxIdleTime(c.Replica.MaxIdleTime)
		}
		err = db.Use(pluginDbResolver)
		if err != nil {
			return nil, err
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if c.Source.MaxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(c.Source.MaxIdleConn)
	}
	if c.Source.MaxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(c.Source.MaxOpenConn)
	}
	if c.Source.MaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(c.Source.MaxLifetime)
	}
	if c.Source.MaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(c.Source.MaxIdleTime)
	}
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// SetDBLogger set db logger
func SetDBLogger(db *gorm.DB, l logger.Interface) {
	db.Logger = l
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func newDialector(dialect string, c *Node) (gorm.Dialector, error) {
	var dialector gorm.Dialector

	switch dialect {
	case "mysql":
		dialector = mysql.New(mysql.Config{
			DSN: c.Dsn,
			// DefaultStringSize:         256,   // string 类型字段的默认长度
			// DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			// DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			// DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			// SkipInitializeWithVersion: false, // 根据版本自动配置
		})
	case "postgres":
		dialector = postgres.New(postgres.Config{DSN: c.Dsn})
	case "sqlite3":
		dsn := c.Dsn
		// 路径是否存在
		_, err := os.Stat(dsn)
		if !(err == nil || os.IsExist(err)) {
			if err := os.MkdirAll(path.Dir(dsn), os.ModePerm); err != nil {
				return nil, fmt.Errorf("database mkdir (%s), %+v", dsn, err)
			}
			if _, err := os.Create(dsn); err != nil {
				return nil, fmt.Errorf("database create DB(%s), %+v", dsn, err)
			}
		}
		dialector = NewSqlite3(dsn)
	default:
		return nil, errors.New("please select database driver one of [mysql|postgres|sqlite3|custom], if use sqlite3, build tags with mysql|postgres|sqlite3!")
	}
	return dialector, nil
}
