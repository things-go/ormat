//go:build !nomssql
// +build !nomssql

package database

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func newMssql(c *Config) gorm.Dialector {
	return sqlserver.Open(c.Dsn)
}
