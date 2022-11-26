package runtime

import (
	"errors"
	stdlog "log"
	"os"
	"time"

	"github.com/google/wire"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/pkg/database"
	"github.com/things-go/ormat/pkg/deploy"
)

var DbSet = wire.NewSet(NewDb, NewDbConfig)

func NewDb(dbCfg *config.Database) (*gorm.DB, error) {
	if dbCfg == nil {
		return nil, errors.New("database config must be set")
	}
	gc := &gorm.Config{}
	if !deploy.IsRelease() {
		gc.Logger = logger.New(stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		})
	}

	return database.New(database.Config{
		Dialect: dbCfg.Dialect,
		Dsn:     dbCfg.Dsn(),
	}, gc)
}

func NewDbConfig(c *config.Config) *config.Database { return c.Database }
