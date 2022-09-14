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

const protobufTpl = `
/*
{{- $tableName := .TableName}}
// {{.StructName}} {{.StructComment}}
message {{.StructName}} { 
{{- range $index, $field := .Fields}}
    {{- if $field.ColumnComment}} 
	// {{$field.ColumnComment}} 
	{{- end}}
	{{$field.ColumnDataType}} {{$field.ColumnName}} = {{$index}};
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
*/
`

// {{$combine := (printf "%s_%s" $field.TableName $field.ColumnName)}}
// {{ $combine }}
var TableNameTpl = template.Must(template.New("tableNameTpl").Parse(tableNameTpl))
var ColumnNameTpl = template.Must(template.New("columnNameTpl").Parse(columnNameTpl))
var ProtobufTpl = template.Must(template.New("protobufTpl").Parse(protobufTpl))
