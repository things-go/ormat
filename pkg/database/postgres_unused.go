//go:build nopostgres
// +build nopostgres

package database

import (
	"gorm.io/gorm"
)

func newPostgres(*Config) gorm.Dialector {
	panic("please build tags without nopostgres!")
}
