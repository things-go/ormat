//go:build sqlite3
// +build sqlite3

package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newSqlite3(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}
