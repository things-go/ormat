// {{.StructName}} {{.Comment}}
type {{.StructName}} struct {
    {{- range $field := .StructFields}}
        {{$field.FieldName}} {{$field.FieldType}}
    {{- end}}
}