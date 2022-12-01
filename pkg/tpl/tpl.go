package tpl

import (
	"embed"
	"text/template"
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

var Template = template.Must(template.New("xx").
	Funcs(template.FuncMap{"add": func(a, b int) int { return a + b }}).
	ParseFS(Static, "template/layout/*"))
