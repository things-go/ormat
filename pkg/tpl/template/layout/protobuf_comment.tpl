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
