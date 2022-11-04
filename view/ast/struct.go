package ast

import (
	"encoding/json"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cast"
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
	StructName    string          // 结构体名
	StructComment string          // 结构体注释
	TableName     string          // 表名
	AbbrTableName string          // 表缩写表名
	Fields        []ProtobufField // 字段列表
	Enums         []ProtobufEnum  // 枚举列表
}

type ProtobufField struct {
	ColumnComment  string // 列注释
	ColumnDataType string // 列数据类型
	ColumnName     string // 列名称
	IsTimestamp    bool   // 是否是时间类型
	Annotation     string // 注解
}

type ProtobufEnum struct {
	EnumName    string              // 枚举名称 表名+列名
	EnumComment string              // 注释
	EnumFields  []ProtobufEnumField // 枚举字段
}

type ProtobufEnumField struct {
	Id      int    // 段序号
	Name    string // 段名称 uppercase(表名_列名_段名)
	Comment string // 段注释
}

type ProtobufEnumFieldSlice []ProtobufEnumField

func (p ProtobufEnumFieldSlice) Len() int           { return len(p) }
func (p ProtobufEnumFieldSlice) Less(i, j int) bool { return p[i].Id < p[j].Id }
func (p ProtobufEnumFieldSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (s *Struct) BuildProtobufTemple() string {
	var buf strings.Builder

	_ = HelperTpl.Execute(&buf, s.intoProtobufMessage())
	return buf.String()
}

func (s *Struct) intoProtobufMessage() *ProtobufMessage {
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
		Fields:        make([]ProtobufField, 0, len(s.StructFields)),
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

		protobufEnum := parseEnumComment(s.StructName, s.TableName, field.ColumnName, field.FieldComment)
		if protobufEnum != nil {
			pm.Enums = append(pm.Enums, *protobufEnum)
		}
	}
	return pm
}

// t.Logf("%#v", rEnum.FindStringSubmatch(` 11 [@enum:{"0":["none"],"1":["expenditure","支出"],"2":["income","收入"]}] 11k l23123 人11`))
var rEnum = regexp.MustCompile(`^.*?\[@.*?(?i:(?:enum|status)+):\s*(.*)\].*?`)

func parseEnumComment(structName, tableName, columnName, comment string) *ProtobufEnum {
	enumCommentString := func(comment string) string {
		match := rEnum.FindStringSubmatch(comment)
		if len(match) == 2 {
			return strings.TrimSpace(match[1])
		}
		return ""
	}

	str := enumCommentString(comment)
	if str == "" {
		return nil
	}
	var mp map[string][]string

	err := json.Unmarshal([]byte(str), &mp)
	if err != nil {
		return nil
	}
	if len(mp) == 0 {
		return nil
	}
	protobufEnum := ProtobufEnum{
		EnumName:    structName + columnName,
		EnumComment: comment,
		EnumFields:  make([]ProtobufEnumField, 0, len(mp)),
	}
	for k, v := range mp {
		protobufEnumField := ProtobufEnumField{
			Id:      cast.ToInt(k),
			Name:    "",
			Comment: "",
		}
		if len(v) > 0 {
			protobufEnumField.Name = strings.ToUpper(tableName + "_" + columnName + "_" + strings.ReplaceAll(v[0], " ", "_"))
		}
		if len(v) > 1 {
			protobufEnumField.Comment = v[1]
		}
		protobufEnum.EnumFields = append(protobufEnum.EnumFields, protobufEnumField)
	}
	sort.Sort(ProtobufEnumFieldSlice(protobufEnum.EnumFields))
	return &protobufEnum
}
