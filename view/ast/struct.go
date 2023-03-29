package ast

import (
	"fmt"
	"strings"
)

// Struct define a struct
type Struct struct {
	StructName         string                 // struct name
	StructComment      string                 // struct comment
	StructFields       []Field                // struct field list
	TableName          string                 // struct table name in database.
	AbbrTableName      string                 // struct abbreviate table name
	CreateTableSQL     string                 // create table SQL
	ProtoMessageFields []ProtobufMessageField // proto message field
}

func ParseProtobuf(structFields []Field, enableGogo, enableSea bool) []ProtobufMessageField {
	// 转成注解
	intoAnnotation := func(annotations []string) string {
		annotation := ""
		if len(annotations) > 0 {
			annotation = "[" + strings.Join(annotations, ", ") + "]"
		}
		return annotation
	}

	protobufMessageFields := make([]ProtobufMessageField, 0, len(structFields))
	tmpAnnotations := make([]string, 0, 16)
	for _, field := range structFields {
		tmpAnnotations = tmpAnnotations[:0]
		dataType := field.ColumnGoType
		// 转换成 proto 类型
		switch dataType {
		case "time.Time":
			if enableGogo {
				dataType = "google.protobuf.Timestamp"
				tmpAnnotations = append(tmpAnnotations, `(gogoproto.stdtime) = true`, `(gogoproto.nullable) = false`)
			} else {
				dataType = "int64"
				tmpAnnotations = append(tmpAnnotations, `(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { type: [ INTEGER ] }`)
			}
			if enableSea {
				tmpAnnotations = append(tmpAnnotations, fmt.Sprintf(`(things_go.seaql.field) = { type: "%s" }`, field.Type))
			}
			protobufMessageFields = append(protobufMessageFields, ProtobufMessageField{
				FieldDataType:   dataType,
				FieldName:       field.ColumnName,
				FieldComment:    field.FieldComment,
				FieldAnnotation: intoAnnotation(tmpAnnotations),
				IsTimestamp:     !enableGogo,
			})
			continue
		case "uint16", "uint8", "uint":
			dataType = "uint32"
		case "int16", "int8", "int":
			dataType = "int32"
		case "float64":
			dataType = "double"
		case "float32":
			dataType = "float"
		case "[]byte":
			dataType = "bytes"
		case "int64", "uint64":
			tmpAnnotations = append(tmpAnnotations,
				`(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { type: [ INTEGER ] }`)
		}
		if enableSea {
			tmpAnnotations = append(tmpAnnotations, fmt.Sprintf(`(things_go.seaql.field) = { type: "%s" }`, field.Type))
		}

		protobufMessageFields = append(protobufMessageFields, ProtobufMessageField{
			FieldDataType:   dataType,
			FieldName:       field.ColumnName,
			FieldComment:    field.FieldComment,
			FieldAnnotation: intoAnnotation(tmpAnnotations),
			IsTimestamp:     false,
		})
	}
	return protobufMessageFields
}
