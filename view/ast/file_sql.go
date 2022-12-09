package ast

type TableAttribute struct {
	Name           string // table name, 表名
	Comment        string // table comment, 表注释
	CreateTableSQL string // Create SQL statements, 创建表的sql语句
}

type SqlFile struct {
	Version string
	Tables  []TableAttribute
}
