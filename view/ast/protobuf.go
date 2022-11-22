package ast

import (
	"encoding/json"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cast"
)

type ProtobufEnumField struct {
	Id      int    // 段序号
	Name    string // 段名称 uppercase(表名_列名_段名)
	Mapping string // 段映射值
	Comment string // 段注释
}

type ProtobufEnum struct {
	EnumName    string              // 枚举名称 表名+列名
	EnumComment string              // 注释
	EnumFields  []ProtobufEnumField // 枚举字段
}

type ProtobufMessageField struct {
	FieldDataType   string // 列数据类型
	FieldName       string // 列名称
	FieldComment    string // 列注释
	FieldAnnotation string // 列注解
	IsTimestamp     bool   // 是否是时间类型
}

type ProtobufMessage struct {
	StructName    string                 // 结构体名
	StructComment string                 // 结构体注释
	TableName     string                 // 表名
	AbbrTableName string                 // 表名缩写
	Fields        []ProtobufMessageField // 字段列表
	Enums         []*ProtobufEnum        // 枚举列表(解析注释中)
}

type ProtobufEnumFieldSlice []ProtobufEnumField

func (p ProtobufEnumFieldSlice) Len() int           { return len(p) }
func (p ProtobufEnumFieldSlice) Less(i, j int) bool { return p[i].Id < p[j].Id }
func (p ProtobufEnumFieldSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// t.Logf("%#v", rEnum.FindStringSubmatch(` 11 [@enum:{"0":["none"],"1":["expenditure","支出"],"2":["income","收入"]}] 11k l23123 人11`))
var rEnum = regexp.MustCompile(`^.*?\[@.*?(?i:(?:enum|status)+):\s*(.*)\].*?`)

func parseEnumComment(structName, tableName, fieldName, columnName, comment string) *ProtobufEnum {
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
		EnumName:    structName + fieldName,
		EnumComment: comment,
		EnumFields:  make([]ProtobufEnumField, 0, len(mp)),
	}
	for k, v := range mp {
		protobufEnumField := ProtobufEnumField{
			Id:      cast.ToInt(k),
			Name:    "",
			Mapping: "",
			Comment: "",
		}
		if len(v) > 0 {
			protobufEnumField.Name = strings.ToUpper(tableName + "_" + columnName + "_" + strings.ReplaceAll(v[0], " ", "_"))
		}
		if len(v) > 1 {
			protobufEnumField.Mapping = v[1]
			protobufEnumField.Comment = v[1]
		}
		if len(v) > 2 && v[2] != "" {
			if protobufEnumField.Comment != "" {
				protobufEnumField.Comment = protobufEnumField.Comment + "," + v[2]
			} else {
				protobufEnumField.Comment = v[2]
			}
		}
		protobufEnum.EnumFields = append(protobufEnum.EnumFields, protobufEnumField)
	}
	sort.Sort(ProtobufEnumFieldSlice(protobufEnum.EnumFields))
	return &protobufEnum
}
