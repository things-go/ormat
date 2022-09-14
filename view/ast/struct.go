package ast

import (
	"strings"
)

// Struct define a struct
type Struct struct {
	StructName     string  // struct name
	StructComment  string  // struct comment
	StructFields   []Field // struct field list
	TableName      string  // table name
	CreateTableSQL string  // create table SQL
}

// AddStructFields Add one or more fields
func (s *Struct) AddStructFields(e ...Field) *Struct {
	s.StructFields = append(s.StructFields, e...)
	return s
}

func (s *Struct) BuildTableNameTemplate() string {
	type tpl struct {
		TableName  string
		StructName string
	}

	var buf strings.Builder

	_ = TableNameTpl.Execute(&buf, tpl{
		TableName:  s.TableName,
		StructName: s.StructName,
	})
	return buf.String()
}

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
	buf.WriteString("// " + s.StructName + " " + comment + delimLF)
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

type ProtobufMessage struct {
	StructName    string
	StructComment string
	TableName     string
	AbbrTableName string
	Fields        []ProtobufField
}

type ProtobufField struct {
	ColumnComment  string
	ColumnDataType string
	ColumnName     string
	IsTimestamp    bool
	Annotation     string
}

func (s *Struct) BuildProtobufTemple() string {
	var buf strings.Builder

	_ = HelperTpl.Execute(&buf, s.intoProtobufMessage())
	return buf.String()
}

func (s *Struct) intoProtobufMessage() *ProtobufMessage {
	pm := &ProtobufMessage{
		StructName:    s.StructName,
		StructComment: s.StructComment,
		TableName:     s.TableName,
		AbbrTableName: intoAbbrTableName(s.TableName),
		Fields:        make([]ProtobufField, 0, len(s.StructFields)),
	}

	intoAnnotation := func(annotations []string) string {
		annotation := ""
		if len(annotations) > 0 {
			annotation = "[" + strings.Join(annotations, ", ") + "]"
		}
		return annotation
	}

	for _, field := range s.StructFields {
		var tmpAnnotations []string
		dataType := field.ColumnDataType
		switch dataType {
		case "time.Time":
			dataType = "google.protobuf.Timestamp"
			pm.Fields = append(pm.Fields,
				ProtobufField{
					ColumnComment:  field.FieldComment,
					ColumnDataType: "google.protobuf.Timestamp",
					ColumnName:     field.ColumnName,
					IsTimestamp:    false,
					Annotation:     intoAnnotation([]string{`(gogoproto.stdtime) = true`, `(gogoproto.nullable) = false`}),
				},
				ProtobufField{
					ColumnComment:  field.FieldComment,
					ColumnDataType: "int64",
					ColumnName:     field.ColumnName,
					IsTimestamp:    true,
					Annotation:     intoAnnotation([]string{`(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { type: [ INTEGER ] }`}),
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

		pm.Fields = append(pm.Fields, ProtobufField{
			ColumnComment:  field.FieldComment,
			ColumnDataType: dataType,
			ColumnName:     field.ColumnName,
			IsTimestamp:    false,
			Annotation:     intoAnnotation(tmpAnnotations),
		})
	}
	return pm
}

func intoAbbrTableName(tableName string) string {
	ss := strings.Split(tableName, "_")
	tableName = ""
	for _, vv := range ss {
		if len(vv) > 0 {
			tableName += string(vv[0])
		}
	}
	return tableName
}
