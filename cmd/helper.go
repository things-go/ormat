package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"github.com/things-go/ormat/pkg/tpl"
	"github.com/things-go/ormat/view"
	"github.com/things-go/ormat/view/driver"
)

var validate = validator.New()

func init() {
	validate.SetTagName("binding")
}

type DbConfig struct {
	DB               *gorm.DB
	Dialect          string
	DbName           string
	TableNames       []string
	CustomDefineType map[string]string
	view.Config
}

func NewFromDatabase(c *DbConfig) (*view.View, error) {
	var m view.DBModel
	switch c.Dialect {
	case "mysql":
		m = &driver.MySQL{
			DB:               c.DB,
			DbName:           c.DbName,
			TableNames:       c.TableNames,
			CustomDefineType: c.CustomDefineType,
		}
	case "sqlite3":
		m = &driver.SQLite{
			DB:               c.DB,
			DbName:           c.DbName,
			TableNames:       c.TableNames,
			CustomDefineType: c.CustomDefineType,
		}
	default:
		return nil, errors.New("database not found, please check database.dialect (mysql, sqlite3, mssql)")
	}
	return view.New(m, c.Config), nil
}

func intoFilename(dir, filename, suffix string) string {
	suffix = strings.TrimSpace(suffix)
	if suffix != "" && !strings.HasPrefix(suffix, ".") {
		suffix = "." + suffix
	}
	return filepath.Join(dir, filename) + suffix
}

func getModelTemplate(filename, suffix string) (*tpl.TemplateMapping, error) {
	return getMappingTemplate(tpl.BuiltInModelMapping, filename, suffix)
}

func getEnumTemplate(filename, suffix string) (*tpl.TemplateMapping, error) {
	return getMappingTemplate(tpl.BuiltInEnumMapping, filename, suffix)
}
func getMappingTemplate(mapping map[string]tpl.TemplateMapping, filename, suffix string) (*tpl.TemplateMapping, error) {
	if t, ok := mapping[filename]; ok {
		if suffix != "" {
			t.Suffix = suffix
		}
		return &t, nil
	}
	t, err := tpl.ParseTemplateFromFile(filename)
	if err != nil {
		return nil, err
	}
	return &tpl.TemplateMapping{
		Template: t,
		Suffix:   suffix,
	}, nil
}

func JSON(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	quick.Highlight(os.Stdout, string(b), "JSON", "terminal", "solarized-dark") // nolint
}
