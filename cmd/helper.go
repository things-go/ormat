package cmd

import (
	"errors"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/go-playground/validator/v10"
	"github.com/things-go/log"

	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/pkg/deploy"
	"github.com/things-go/ormat/pkg/tpl"
	"github.com/things-go/ormat/runtime"
	"github.com/things-go/ormat/view"
	"github.com/things-go/ormat/view/driver"
)

var validate = validator.New()

func init() {
	validate.SetTagName("binding")
}

func GetViewModel(rt *runtime.Runtime) view.DBModel {
	c := config.Global
	dbCfg := c.Database
	switch dbCfg.Dialect {
	case "mysql":
		return &driver.MySQL{
			DB:               rt.DB,
			DbName:           dbCfg.DbName(),
			TableNames:       c.TableNames,
			CustomDefineType: c.TypeDefine,
		}
	case "sqlite3":
		return &driver.SQLite{
			DB:               rt.DB,
			DbName:           dbCfg.DbName(),
			TableNames:       c.TableNames,
			CustomDefineType: c.TypeDefine,
		}
	default:
		panic("database not found, please check database.dialect (mysql, sqlite3, mssql)")
	}
}

func setupBase(c *config.Config) {
	deploy.MustSetDeploy(c.Deploy)
	log.ReplaceGlobals(log.NewLogger(log.WithConfig(log.Config{
		Level:  "info",
		Format: "console",
	})))
}

func GetFilenameSuffix(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ".proto"
	}
	if !strings.HasPrefix(s, ".") {
		s = "." + s
	}
	return s
}

func intoFilename(dir, filename, suffix string) string {
	return filepath.Join(dir, filename) + GetFilenameSuffix(suffix)
}

func parseTemplateFromFile(filename string) (*template.Template, error) {
	if filename == "" {
		return nil, errors.New("not found template file")
	}
	tt, err := template.New("custom").
		Funcs(tpl.TemplateFuncs).
		ParseFiles(filename)
	if err != nil {
		return nil, err
	}
	ts := tt.Templates()
	if len(ts) == 0 {
		return nil, errors.New("not found any template")
	}
	return ts[0], nil
}
