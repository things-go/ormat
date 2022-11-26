package config

import (
	"errors"
	"fmt"
	"path/filepath"
)

// Database connect information
type Database struct {
	Dialect  string `yaml:"dialect" json:"dialect" binding:"oneof=mysql sqlite3"` // mysql, sqlite3
	Host     string `yaml:"host" json:"host" binding:"required"`                  // Host. 地址
	Port     int    `yaml:"port" json:"port" binding:"required"`                  // Port 端口号
	Username string `yaml:"username" json:"username" binding:"required"`          // Username 用户名
	Password string `yaml:"password" json:"password" binding:"required"`          // Password 密码
	Db       string `yaml:"db" json:"db" binding:"required"`                      // Database 数据库名
	Options  string `yaml:"options" json:"options"`                               // Options ?号后面, 如果为空, 则为 charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=True

	// primary field
	dsn    string // dsn
	dbName string // real db name
}

func (c *Database) Dsn() string    { return c.dsn }
func (c *Database) DbName() string { return c.dbName }
func (c *Database) Parse() error {
	switch c.Dialect {
	case "mysql":
		if c.Options == "" {
			c.Options = "charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=True"
		}
		c.dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
			c.Username, c.Password,
			c.Host, c.Port,
			c.Db, c.Options,
		)
		c.dbName = c.Db
		return nil
	case "sqlite3":
		dbName := filepath.Base(c.Db)
		if dbName == "" {
			return errors.New("empty sqlite3 db name")
		}
		c.dsn = c.Db
		c.dbName = dbName
		return nil
	default:
		return errors.New("database not found, please check database.dialect (mysql, sqlite3, mssql)")
	}
}
