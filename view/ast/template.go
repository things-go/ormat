package ast

import "html/template"

// interval
const delimTab = "\t"
const delimLF = "\n"

const tableNameTpl = `
// TableName implement schema.Tabler interface
func (*{{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
`

const columnNameTpl = `
// {{.StructName}}Columns get sql column name
var {{.StructName}}Columns = struct { 
{{- range $field := .Fields}}
	{{$field.Name}} string
{{- end}}    
}{
{{- range $field := .Fields}}
	{{$field.Name}}:"{{$field.ColumnName}}",  
{{- end}}           
	}
`

var TableNameTpl = template.Must(template.New("tableNameTpl").Parse(tableNameTpl))
var ColumnNameTpl = template.Must(template.New("columnNameTpl").Parse(columnNameTpl))
