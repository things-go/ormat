// Code generated by ormat. DO NOT EDIT.
// version: {{.Version}}

package {{.PackageName}}

{{if .Imports}}
import (
{{- range $k, $v := .Imports}}
    {{$k}}
{{- end}}
)
{{end}}

{{- $hasColumn := .HasColumn}}

{{- range $e := .Structs}}
// {{$e.StructName}} {{$e.StructComment}}
type {{$e.StructName}} struct {
{{- range $field := $e.StructFields}}
    {{$field.FieldName}} {{$field.FieldType}} {{if $field.FieldTag}}`{{$field.FieldTag}}`{{end}} {{if $field.FieldComment}}// {{$field.FieldComment}}{{end}}
{{- end}}
}

// TableName implement schema.Tabler interface
func (*{{$e.StructName}}) TableName() string {
	return "{{$e.TableName}}"
}

{{- $tableName := $e.TableName}}
{{- if $hasColumn}}
// Select{{$e.StructName}} database column name.
var Select{{$e.StructName}} = []string {
{{- range $field := $e.StructFields}}
	{{- if $field.IsTimestamp}}
	{{- if $field.IsNullable}}
	{{if $field.IsSkipColumn}}// {{end}}"IFNULL(UNIX_TIMESTAMP(`{{$tableName}}`.`{{$field.ColumnName}}`), 0) AS `{{$field.ColumnName}}`",
	{{- else}}
	{{if $field.IsSkipColumn}}// {{end}}"UNIX_TIMESTAMP(`{{$tableName}}`.`{{$field.ColumnName}}`) AS `{{$field.ColumnName}}`",
	{{- end}}
	{{- else}}
	{{if $field.IsSkipColumn}}// {{end}}"`{{$tableName}}`.`{{$field.ColumnName}}`",
	{{- end}}
{{- end}}
}
{{- end}}
{{- end}}

