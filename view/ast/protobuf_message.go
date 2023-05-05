package ast

type ProtobufMessageField struct {
	FieldDataType   string // 列数据类型
	FieldName       string // 列名称
	FieldComment    string // 列注释
	FieldAnnotation string // 列注解
}
