package tpl

import (
	"embed"
	"errors"
	"text/template"

	"github.com/things-go/ormat/pkg/utils"
)

//go:embed template
var Static embed.FS

var (
	TemplateFuncs = template.FuncMap{
		"add":            func(a, b int) int { return a + b },
		"snakecase":      func(s string) string { return utils.SnakeCase(s, false) },
		"kebabcase":      func(s string) string { return utils.Kebab(s, false) },
		"camelcase":      func(s string) string { return utils.CamelCase(s, false) },
		"smallcamelcase": func(s string) string { return utils.SmallCamelCase(s, false) },
	}

	Template = template.Must(
		template.New("components").
			Funcs(TemplateFuncs).
			ParseFS(Static, "template/layout/*"),
	)

	Model  = Template.Lookup("model.tpl")
	SqlDDL = Template.Lookup("sql_ddl.tpl")
)

type TemplateMapping struct {
	Template *template.Template
	Suffix   string
}

var BuiltInModelMapping = map[string]TemplateMapping{
	"__in_go": {Model, ".go"},
}

func ParseTemplateFromFile(filename string) (*template.Template, error) {
	if filename == "" {
		return nil, errors.New("required template filename")
	}
	tt, err := template.New("custom").
		Funcs(TemplateFuncs).
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
