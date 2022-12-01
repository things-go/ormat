package ast

import (
	"strings"
)

// Struct define a struct
type Struct struct {
	StructName     string           // struct name
	StructComment  string           // struct comment
	StructFields   []Field          // struct field list
	TableName      string           // struct table name in database.
	CreateTableSQL string           // create table SQL
	protoMessage   *ProtobufMessage // proto message.
}

// AddStructFields Add one or more fields
func (s *Struct) AddStructFields(e ...Field) *Struct {
	s.StructFields = append(s.StructFields, e...)
	return s
}

// Build Get the struct data.
func (s *Struct) Build() string {
	buf := &strings.Builder{}

	comment := s.StructComment
	if comment != "" {
		comment = strings.ReplaceAll(strings.TrimSpace(comment), "\n", "\r\n// ")
	} else {
		comment = "..."
	}
	// comment
	buf.WriteString("// " + s.StructName + delimSpace + comment + delimLF)
	buf.WriteString("type\t" + s.StructName + "\tstruct {" + delimLF)

	// field every line
	mp := make(map[string]struct{}, len(s.StructFields))
	for _, field := range s.StructFields {
		if _, ok := mp[field.FieldName]; !ok {
			mp[field.FieldName] = struct{}{}
			buf.WriteString("\t\t" + field.BuildLine() + delimLF)
		}
	}
	buf.WriteString("}" + delimLF)

	return buf.String()
}

// BuildTableNameTemplate struct implement schema.Tabler.
func (s *Struct) BuildTableNameTemplate() string {
	type tmpTpl struct {
		TableName  string
		StructName string
	}

	var buf strings.Builder

	_ = TableNameTpl.Execute(&buf, tmpTpl{
		TableName:  s.TableName,
		StructName: s.StructName,
	})
	return buf.String()
}

// BuildColumnNameTemplate field name mapping column name which in database.
func (s *Struct) BuildColumnNameTemplate() string {
	type tpl struct {
		StructName string
		Fields     []Field
	}
	var buf strings.Builder

	_ = ColumnNameTpl.Execute(&buf, &tpl{
		StructName: s.StructName,
		Fields:     s.StructFields,
	})
	return buf.String()
}

func (s *Struct) BuildSQL() string {
	buf := &strings.Builder{}

	comment := s.StructComment
	if comment != "" {
		comment = strings.ReplaceAll(strings.TrimSpace(comment), "\n", "\r\n// ")
	} else {
		comment = "..."
	}
	// comment
	buf.WriteString("-- " + s.StructName + " " + comment + delimLF)
	// sql
	buf.WriteString(s.CreateTableSQL + ";" + delimLF)
	return buf.String()
}

func (s *Struct) BuildProtobufTemplate() string {
	var buf strings.Builder

	s.parseProtobufMessage()
	_ = ProtobufTpl.Execute(&buf, s.protoMessage)
	buf.WriteString(s.buildProtobufEnumTemplate(true))
	return buf.String()
}

func (s *Struct) BuildProtobufEnumTemplate() string {
	return s.buildProtobufEnumTemplate(false)
}

func (s *Struct) buildProtobufEnumTemplate(isAnnotation bool) string {
	type tpl struct {
		Enums        []*ProtobufEnum
		IsAnnotation bool
		Package      string
		Options      map[string]string
	}
	var buf strings.Builder

	s.parseProtobufMessage()
	if len(s.protoMessage.Enums) > 0 {
		_ = ProtobufEnumTpl.Execute(&buf, &tpl{
			Enums:        s.protoMessage.Enums,
			IsAnnotation: isAnnotation,
		})
	}
	return buf.String()
}

func (s *Struct) BuildProtobufEnumMappingTemplate() string {
	type tpl struct {
		Enums []*ProtobufEnum
	}
	var buf strings.Builder

	s.parseProtobufMessage()
	if len(s.protoMessage.Enums) > 0 {
		_ = ProtobufEnumMappingTpl.Execute(&buf, &tpl{
			Enums: s.protoMessage.Enums,
		})
	}
	return buf.String()
}

func (s *Struct) parseProtobufMessage() {
	if s.protoMessage != nil {
		return
	}

	// 获取表名缩写
	intoAbbrTableName := func(tableName string) string {
		ss := strings.Split(tableName, "_")
		tableName = ""
		for _, vv := range ss {
			if len(vv) > 0 {
				tableName += string(vv[0])
			}
		}
		return tableName
	}
	// 转获成注解
	intoAnnotation := func(annotations []string) string {
		annotation := ""
		if len(annotations) > 0 {
			annotation = "[" + strings.Join(annotations, ", ") + "]"
		}
		return annotation
	}

	pm := &ProtobufMessage{
		StructName:    s.StructName,
		StructComment: s.StructComment,
		TableName:     s.TableName,
		AbbrTableName: intoAbbrTableName(s.TableName),
		Fields:        make([]ProtobufMessageField, 0, len(s.StructFields)),
		Enums:         make([]*ProtobufEnum, 0, 32),
	}
	for _, field := range s.StructFields {
		var tmpAnnotations []string
		dataType := field.ColumnDataType
		// 转换成 proto 类型
		switch dataType {
		case "time.Time":
			dataType = "google.protobuf.Timestamp"
			pm.Fields = append(pm.Fields,
				ProtobufMessageField{
					FieldDataType:   dataType,
					FieldName:       field.ColumnName,
					FieldComment:    field.FieldComment,
					FieldAnnotation: intoAnnotation([]string{`(gogoproto.stdtime) = true`, `(gogoproto.nullable) = false`}),
					IsTimestamp:     false,
				},
				ProtobufMessageField{
					FieldDataType:   "int64",
					FieldName:       field.ColumnName,
					FieldComment:    field.FieldComment,
					FieldAnnotation: intoAnnotation([]string{`(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { type: [ INTEGER ] }`}),
					IsTimestamp:     true,
				},
			)
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

		pm.Fields = append(pm.Fields, ProtobufMessageField{
			FieldDataType:   dataType,
			FieldName:       field.ColumnName,
			FieldComment:    field.FieldComment,
			FieldAnnotation: intoAnnotation(tmpAnnotations),
			IsTimestamp:     false,
		})

		protobufEnum := parseEnumComment(s.StructName, s.TableName, field.FieldName, field.ColumnName, field.FieldComment)
		if protobufEnum != nil {
			pm.Enums = append(pm.Enums, protobufEnum)
		}
	}
	s.protoMessage = pm
}
