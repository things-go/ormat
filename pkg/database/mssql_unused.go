//go:build nomssql
// +build nomssql

package database

import (
	"gorm.io/gorm"
)

func newMssql(*Config) gorm.Dialector {
	panic("please build tags without nomssql!")
}
