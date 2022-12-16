package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/quick"
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

func setupBase2(dp string) {
	deploy.MustSetDeploy(dp)
	log.ReplaceGlobals(log.NewLogger(log.WithConfig(log.Config{
		Level:  "info",
		Format: "console",
	})))
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
