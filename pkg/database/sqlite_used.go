//go:build sqlite3
// +build sqlite3

package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqlite3(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}
