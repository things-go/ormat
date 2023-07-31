package command

import (
	"errors"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/things-go/ens/driver"
	driverMysql "github.com/things-go/ens/driver/mysql"
)

func NewDriver(URL string) (driver.Driver, error) {
	var m driver.Driver

	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "mysql":
		m = &driverMysql.MySQL{
			URL: URL,
		}

	// case "sqlite3":
	// 	m = &driver.SQLite{
	// 		DB:               c.DB,
	// 		DbName:           c.DbName,
	// 		TableNames:       c.TableNames,
	// 	}
	default:
		return nil, errors.New("database not found, please check database.dialect (mysql, sqlite3, mssql)")
	}
	return m, nil
}

func joinFilename(dir, filename, suffix string) string {
	suffix = strings.TrimSpace(suffix)
	if suffix != "" && !strings.HasPrefix(suffix, ".") {
		suffix = "." + suffix
	}
	return filepath.Join(dir, filename) + suffix
}

// WriteFile writes data to a file named by filename.
// If the file does not exist, WriteFile creates it
// and its upper level paths.
func WriteFile(filename string, data []byte) error {
	if err := os.MkdirAll(path.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0655)
}
