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
