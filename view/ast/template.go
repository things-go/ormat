package ast

import "text/template"

const tableNameTemplate = `
// TableName implement schema.Tabler interface
func (*{{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
`

const columnNameTemplate = `
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
`

const protobufTemplate = `
/* protobuf and gorm field helper
{{- $tableName := .TableName}}
{{- $abbrTableName := .AbbrTableName}}
// {{.StructName}}Columns database column name.
var {{.StructName}}Columns = []string {
{{- range $field := .Fields}}
	{{- if $field.IsTimestamp}}
	"UNIX_TIMESTAMP({{$field.FieldName}}) AS {{$field.FieldName}}",  
	{{- else}}
	"{{$field.FieldName}}",  
	{{- end}}
{{- end}}
}
// {{.StructName}}ColumnsWithTable database column name with table prefix
var {{.StructName}}ColumnsWithTable = []string {
{{- range $field := .Fields}}
	{{- if $field.IsTimestamp}}
	"UNIX_TIMESTAMP({{$tableName}}.{{$field.FieldName}}) AS {{$tableName}}_{{$field.FieldName}}", 
	{{- else}}
	"{{$tableName}}.{{$field.FieldName}} AS {{$tableName}}_{{$field.FieldName}}", 
	{{- end}}
{{- end}}
}
// {{.StructName}}ColumnsWithAbbrTable database column name with abbr table prefix
var {{.StructName}}ColumnsWithAbbrTable = []string {
{{- range $field := .Fields}}
	{{- if $field.IsTimestamp}}
	"UNIX_TIMESTAMP({{$abbrTableName}}.{{$field.FieldName}}) AS {{$abbrTableName}}_{{$field.FieldName}}", 
	{{- else}}
	"{{$abbrTableName}}.{{$field.FieldName}} AS {{$abbrTableName}}_{{$field.FieldName}}",
	{{- end}}
{{- end}}
}

// {{.StructName}} {{.StructComment}}
message {{.StructName}} { 
{{- range $index, $field := .Fields}}
    {{- if $field.FieldComment}} 
	// {{$field.FieldComment}} 
	{{- end}}
	{{$field.FieldDataType}} {{$field.FieldName}} = {{add $index 1}} {{- if $field.FieldAnnotation}} {{$field.FieldAnnotation}} {{- end}};
{{- end}}    
}
// {{.StructName}}WithTable {{.StructComment}}
message {{.StructName}}WithTable { 
{{- range $index, $field := .Fields}}
    {{- if $field.FieldComment}} 
	// {{$field.FieldComment}} 
	{{- end}}
	{{$field.FieldDataType}} {{$tableName}}_{{$field.FieldName}} = {{add $index 1}} {{- if $field.FieldAnnotation}} {{$field.FieldAnnotation}} {{- end}};
{{- end}}    
}
// {{.StructName}}WithAbbrTable {{.StructComment}}
message {{.StructName}}WithAbbrTable { 
{{- range $index, $field := .Fields}}
    {{- if $field.FieldComment}} 
	// {{$field.FieldComment}} 
	{{- end}}
	{{$field.FieldDataType}} {{$abbrTableName}}_{{$field.FieldName}} = {{add $index 1}} {{- if $field.FieldAnnotation}} {{$field.FieldAnnotation}} {{- end}};
{{- end}}    
}
*/
`
const protobufEnumTemplate = `
{{- if .IsAnnotation}}
/*
{{- end}}
{{- range $e := .Enums}}
// {{$e.EnumName}} {{$e.EnumComment}}
enum {{$e.EnumName}} {
{{- range $ee := $e.EnumFields}}
    {{- if $ee.Comment}} 
	// {{$ee.Comment}}
	{{- end}}
	{{$ee.Name}} = {{$ee.Id}};
{{- end}} 
}
{{- end}}
{{- if .IsAnnotation}}
*/
{{- end}}
`

const protobufEnumMappingTemplate = `
{{- range $e := .Enums}}
// __{{$e.EnumName}}Mapping  {{$e.EnumName}} mapping
var __{{$e.EnumName}}Mapping = map[int]string{
{{- range $ee := $e.EnumFields}}
	{{$ee.Id}}: "{{$ee.Mapping}}",
{{- end}} 
}

// Get{{$e.EnumName}}Desc get mapping description
func Get{{$e.EnumName}}Desc(t int) string {
	return __{{$e.EnumName}}Mapping[t]
}
{{- end}}
`

var TableNameTpl = template.Must(template.New("tableNameTemplate").Parse(tableNameTemplate))
var ColumnNameTpl = template.Must(template.New("columnNameTemplate").Parse(columnNameTemplate))
var ProtobufTpl = template.Must(template.New("protobufTemplate").Funcs(template.FuncMap{"add": func(a, b int) int { return a + b }}).Parse(protobufTemplate))
var ProtobufEnumTpl = template.Must(template.New("protobufEnumTemplate").Parse(protobufEnumTemplate))
var ProtobufEnumMappingTpl = template.Must(template.New("protobufEnumMappingTemplate").Parse(protobufEnumMappingTemplate))
