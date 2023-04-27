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
{{- $hasHelper := .HasHelper}}
{{- $hasAssist := .HasAssist}}
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

{{- if $hasAssist}}
type {{$e.StructName}}Impl struct {
	// private fields
	xTableName string 

	ALL assist.Asterisk
{{- range $field := $e.StructFields}}
    {{$field.FieldName}} assist.{{$field.AssistType}}
{{- end}}

}

var xx_{{$e.StructName}} = New_X_{{$e.StructName}}("{{$e.TableName}}")

func X_{{$e.StructName}}() {{$e.StructName}}Impl {
	return xx_{{$e.StructName}}
}

func New_X_{{$e.StructName}}(tableName string) {{$e.StructName}}Impl {
	return {{$e.StructName}}Impl{
		xTableName: tableName,

		ALL:  assist.NewAsterisk(tableName),

	{{range $field := $e.StructFields}}
		{{$field.FieldName}}: assist.New{{$field.AssistType}}(tableName, "{{$field.ColumnName}}"),
	{{- end}}			
	}
}

func (*{{$e.StructName}}Impl) As(alias string) {{$e.StructName}}Impl {
	return New_X_{{$e.StructName}}(alias)
}

func (*{{$e.StructName}}Impl) X_Model() *{{$e.StructName}} {
	return &{{$e.StructName}}{}
}

func (*{{$e.StructName}}Impl) Xc_Model() assist.Condition {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&{{$e.StructName}}{})
	}
}

func (x *{{$e.StructName}}Impl) X_TableName() string {
	return x.xTableName
}

func X_Select{{$e.StructName}}() assist.Condition {
	x := &xx_{{$e.StructName}}
	return assist.Select(
{{- range $field := $e.StructFields}}
	{{- if $field.IsTimestamp}}
	{{if $field.IsSkipColumn}}// {{end}}x.{{$field.FieldName}}.UnixTimestamp(){{- if $field.IsNullable}}.IfNull(0){{- end}}.As("{{$field.ColumnName}}"),
	{{- else}}
	{{if $field.IsSkipColumn}}// {{end}}x.{{$field.FieldName}},
	{{- end}}
{{- end}}
	)
}

func X_Select{{$e.StructName}}WithPrefix(prefix string) assist.Condition {
	if prefix == "" {
		return X_Select{{$e.StructName}}()
	}
	x := &xx_{{$e.StructName}}
	return assist.Select(
{{- range $field := $e.StructFields}}
	{{if $field.IsSkipColumn}}// {{end}}x.{{$field.FieldName}}{{- if $field.IsTimestamp}}.UnixTimestamp(){{- if $field.IsNullable}}.IfNull(0){{- end}}{{- end}}.AsWithPrefix(prefix),
{{- end}}
	)
}

{{- end}}

{{- if $hasHelper}}
/* protobuf field helper
// {{$e.StructName}} {{.StructComment}}
message {{$e.StructName}} {
{{- range $index, $field := $e.ProtoMessageFields}}
  {{- if $field.FieldComment}}
  // {{$field.FieldComment}}
  {{- end}}
  {{$field.FieldDataType}} {{$field.FieldName}} = {{add $index 1}} {{- if $field.FieldAnnotation}} {{$field.FieldAnnotation}} {{- end}};
{{- end}}
}
*/
{{- end}}
{{- end}}

