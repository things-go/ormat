// {{.StructName}}Columns field name mapping column name which in database.
var {{.StructName}}Columns = struct {
{{- range $field := .Fields}}
	{{$field.FieldName}} string
{{- end}}
}{
{{- range $field := .Fields}}
	{{$field.FieldName}}:"{{$field.ColumnName}}",
{{- end}}
}
