//go:build nomysql
// +build nomysql

package database

import (
	"gorm.io/gorm"
)

func newMysql(*Config) gorm.Dialector {
	panic("please build tags without nomysql!")
}
