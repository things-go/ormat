//go:build !sqlite3
// +build !sqlite3

package command

import "gorm.io/gorm"

func NewSqlite3(string) gorm.Dialector {
	panic("please build tags with sqlite3!")
}
