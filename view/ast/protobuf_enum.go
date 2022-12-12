package ast

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/spf13/cast"
	"github.com/things-go/log"

	"github.com/things-go/ormat/pkg/matcher"
)

// ProtobufEnumField protobuf enum field
// enum comment format: {"0":["name","mapping","comment"]}
type ProtobufEnumField struct {
	Id      int    // 段序号
	Name    string // 段名称 uppercase(表名_列名_段名)
	Mapping string // 段映射值
	Comment string // 段注释
}

// ProtobufEnumFieldSlice protobuf enum field slice
type ProtobufEnumFieldSlice []ProtobufEnumField

func (p ProtobufEnumFieldSlice) Len() int           { return len(p) }
func (p ProtobufEnumFieldSlice) Less(i, j int) bool { return p[i].Id < p[j].Id }
func (p ProtobufEnumFieldSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// ProtobufEnum protobuf enum
type ProtobufEnum struct {
	EnumName    string              // 枚举名称,表名+列名
	EnumComment string              // 枚举注释
	EnumFields  []ProtobufEnumField // 枚举字段
}

// ParseEnumComment parse enum comment
func ParseEnumComment(structName, tableName, fieldName, columnName, comment string) *ProtobufEnum {
	annotation := matcher.EnumAnnotation(comment)
	if annotation == "" {
		return nil
	}
	mp, err := ParseEnumAnnotation(annotation)
	if err != nil || len(mp) == 0 {
		log.Warnf("🧐 获取到枚举注解解析失败[表:%s, 列: %s, 注解: %s", tableName, columnName, annotation)
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

// ParseEnumAnnotation 解析枚举注解.
func ParseEnumAnnotation(annotation string) (map[string][]string, error) {
	var mp map[string][]string

	err := json.Unmarshal([]byte(annotation), &mp)
	if err != nil {
		return nil, err
	}
	return mp, nil
}
