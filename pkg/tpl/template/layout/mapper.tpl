// Code generated by ormat. DO NOT EDIT.
// version: {{.Version}}

syntax = "proto3";

package {{.Package}};

{{- if .Options}}
{{- range $k, $v := .Options}}
option {{$k}} = "{{$v}}";
{{- end}}
{{- end}}

import "protoc-gen-openapiv2/options/annotations.proto";
import "protosaber/seaql/seaql.proto";

{{- range $e := .Structs}}
// {{$e.StructName}} {{.StructComment}} field
message {{$e.StructName}} {
  option (things_go.seaql.options) = {
    index: [
{{- $indexlen := len $e.SeaIndexes}}
{{- $indexlen := sub $indexlen 1}}
{{- range $index, $field := $e.SeaIndexes}}
      '{{$field}}'{{- if ne $index $indexlen }},{{- end}}
{{- end}}
    ];
};

{{- range $index, $field := $e.ProtoMessageFields}}
  {{- if $field.FieldComment}}
  // {{$field.FieldComment}}
  {{- end}}
  {{$field.FieldDataType}} {{$field.FieldName}} = {{add $index 1}} {{- if $field.FieldAnnotation}} {{$field.FieldAnnotation}} {{- end}};
{{- end}}
}
{{- end}}

