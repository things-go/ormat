package ast

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
}
