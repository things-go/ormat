package tpl

import (
	"embed"
	"text/template"

	"github.com/things-go/ormat/pkg/utils"
)

const (
	ColumnName          = "column_name.tpl"
	ProtobufComment     = "protobuf_comment.tpl"
	ProtobufEnum        = "protobuf_enum.tpl"
	ProtobufEnumMapping = "protobuf_enum_mapping.tpl"
	TableName           = "table_name.tpl"
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
