//go:build !nopostgres
// +build !nopostgres

package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newPostgres(c *Config) gorm.Dialector {
	return postgres.New(postgres.Config{
		DSN: c.Dsn,
	})
}
