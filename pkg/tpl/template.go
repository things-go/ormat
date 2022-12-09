package tpl

import (
	"embed"
	"text/template"

	"github.com/things-go/ormat/pkg/utils"
)

//go:embed template
var Static embed.FS

var TemplateFuncs = template.FuncMap{
	"add":            func(a, b int) int { return a + b },
	"snakecase":      func(s string) string { return utils.SnakeCase(s, false) },
	"kebabcase":      func(s string) string { return utils.Kebab(s, false) },
	"camelcase":      func(s string) string { return utils.CamelCase(s, false) },
	"smallcamelcase": func(s string) string { return utils.SmallCamelCase(s, false) },
}
var Template = template.Must(
	template.New("components").
		Funcs(TemplateFuncs).
		ParseFS(Static, "template/layout/*"),
)

var ProtobufEnumTpl = Template.Lookup("protobuf_enum.tpl")
var ProtobufEnumMappingTpl = Template.Lookup("protobuf_enum_mapping.tpl")
var SqlDDLTpl = Template.Lookup("sql_ddl.tpl")
var ModelTpl = Template.Lookup("model.tpl")
