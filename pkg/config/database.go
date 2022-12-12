package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-sql-driver/mysql"
)

// Database connect information
type Database struct {
	Dialect string `yaml:"dialect" json:"dialect" binding:"oneof=mysql sqlite3"` // mysql, sqlite3
	DSN     string `yaml:"dsn" json:"dsn"`
	Options string `yaml:"options" json:"options"` // Options ?号后面, 如果为空, 则为 charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=True

	// primary field
	dbName string // real db name
}

func (c *Database) DbName() string { return c.dbName }
func (c *Database) Parse() error {
	switch c.Dialect {
	case "mysql":
		if c.Options == "" {
			c.Options = "charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=True"
		}
		idx := strings.Index(c.DSN, "?")
		if idx == -1 {
			c.DSN = fmt.Sprintf("%s?%s", c.DSN, c.Options)
		}
		cc, err := mysql.ParseDSN(c.DSN)
		if err != nil {
			return err
		}
		c.dbName = cc.DBName
		return nil
	case "sqlite3":
		dbName := filepath.Base(c.DSN)
		if dbName == "" {
			return errors.New("empty sqlite3 db name")
		}
		c.dbName = dbName
		return nil
	default:
		return errors.New("database not found, please check database.dialect (mysql, sqlite3, mssql)")
	}
}
