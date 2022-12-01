{{- range $e := .Enums}}
// __{{$e.EnumName}}Mapping {{$e.EnumName}} mapping
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
