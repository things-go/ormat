package ast

import "text/template"

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
	{{$field.FieldName}} string
{{- end}}    
}{
{{- range $field := .Fields}}
	{{$field.FieldName}}:"{{$field.ColumnName}}",  
{{- end}}           
}
`

const helperTpl = `
/* protobuf and gorm field helper
{{- $tableName := .TableName}}
{{- $abbrTableName := .AbbrTableName}}
// {{.StructName}}Columns get sql column name
var {{.StructName}}Columns = []string {
{{- range $field := .Fields}}
	{{- if $field.IsTimestamp}}
	"UNIX_TIMESTAMP({{$field.ColumnName}}) AS {{$field.ColumnName}}",  
	{{- else}}
	"{{$field.ColumnName}}",  
	{{- end}}
{{- end}}
}
// {{.StructName}}ColumnsWithTable get sql column name with table prefix
var {{.StructName}}ColumnsWithTable = []string {
{{- range $field := .Fields}}
	{{- if $field.IsTimestamp}}
	"UNIX_TIMESTAMP({{$tableName}}.{{$field.ColumnName}}) AS {{$tableName}}_{{$field.ColumnName}}", 
	{{- else}}
	"{{$tableName}}.{{$field.ColumnName}} AS {{$tableName}}_{{$field.ColumnName}}", 
	{{- end}}
{{- end}}
}
// {{.StructName}}ColumnsWithAbbrTable get sql column name with abbr table prefix
var {{.StructName}}ColumnsWithAbbrTable = []string {
{{- range $field := .Fields}}
	{{- if $field.IsTimestamp}}
	"UNIX_TIMESTAMP({{$abbrTableName}}.{{$field.ColumnName}}) AS {{$abbrTableName}}_{{$field.ColumnName}}", 
	{{- else}}
	"{{$abbrTableName}}.{{$field.ColumnName}} AS {{$abbrTableName}}_{{$field.ColumnName}}",
	{{- end}}
{{- end}}
}

// {{.StructName}} {{.StructComment}}
message {{.StructName}} { 
{{- range $index, $field := .Fields}}
    {{- if $field.ColumnComment}} 
	// {{$field.ColumnComment}} 
	{{- end}}
	{{$field.ColumnDataType}} {{$field.ColumnName}} = {{$index}} {{- if $field.Annotation}} {{$field.Annotation}} {{- end}};
{{- end}}    
}
// {{.StructName}}WithTable {{.StructComment}}
message {{.StructName}}WithTable { 
{{- range $index, $field := .Fields}}
    {{- if $field.ColumnComment}} 
	// {{$field.ColumnComment}} 
	{{- end}}
	{{$field.ColumnDataType}} {{$tableName}}_{{$field.ColumnName}} = {{$index}} {{- if $field.Annotation}} {{$field.Annotation}} {{- end}};
{{- end}}    
}
// {{.StructName}}WithAbbrTable {{.StructComment}}
message {{.StructName}}WithAbbrTable { 
{{- range $index, $field := .Fields}}
    {{- if $field.ColumnComment}} 
	// {{$field.ColumnComment}} 
	{{- end}}
	{{$field.ColumnDataType}} {{$abbrTableName}}_{{$field.ColumnName}} = {{$index}} {{- if $field.Annotation}} {{$field.Annotation}} {{- end}};
{{- end}}    
}
*/
`

var TableNameTpl = template.Must(template.New("tableNameTpl").Parse(tableNameTpl))
var ColumnNameTpl = template.Must(template.New("columnNameTpl").Parse(columnNameTpl))
var HelperTpl = template.Must(template.New("helperTpl").Parse(helperTpl))
