package view

import (
	"regexp"
	"strings"
)

// ColumnKeyType column key type
type ColumnKeyType int

const (
	ColumnKeyNone            ColumnKeyType = iota // default
	ColumnKeyTypePrimary                          // primary key
	ColumnKeyTypeNormalIndex                      // normal index key
	ColumnKeyTypeUniqueKey                        // unique key
)

// Database database information
type Database struct {
	Name   string  // database name, 数据库名
	Tables []Table // table information, 表信息
}

// TableAttribute database table name, comment and create table sql
type TableAttribute struct {
	Name           string // table name, 表名
	Comment        string // table comment, 表注释
	CreateTableSQL string // Create SQL statements, 创建表的sql语句
}

// Table database table information
type Table struct {
	TableAttribute
	Columns []Column // column information
}

type Tables []Table

func (t Tables) Len() int {
	return len(t)
}

func (t Tables) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

func (t Tables) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Column column information
type Column struct {
	Name            string       // column name
	OrdinalPosition int          // column ordinal position
	DataType        string       // column data type(string,int...)
	ColumnType      string       // column type(varchar(256)...)
	IsNullable      bool         // column is null or not
	IsAutoIncrement bool         // column auto increment or not
	Default         *string      // default value
	Comment         string       // column comment
	Index           []Index      // index list
	ForeignKeys     []ForeignKey // Foreign key list
}

type Columns []Column

func (t Columns) Len() int {
	return len(t)
}

func (t Columns) Less(i, j int) bool {
	return t[i].OrdinalPosition < t[j].OrdinalPosition
}

func (t Columns) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Index database index/unique_index list
type Index struct {
	KeyType    ColumnKeyType // key type
	KeyName    string        // index key name, 索引名称
	IsMulti    bool          // Multiple key, 是否为复合索引
	SeqInIndex int           // union index sequence in index, 复合索引中的序列
	IndexType  string        // index type, 索引类型(比如: BTREE, HASH, FULLTEXT)
}

// ForeignKey Foreign key
type ForeignKey struct {
	TableName  string // Affected tables .
	ColumnName string // Which column of the affected table
}

var nullToPointer = map[string]string{
	"bool":      "*bool",
	"int8":      "*int8",
	"uint8":     "*uint8",
	"int16":     "*int16",
	"uint16":    "*uint16",
	"int32":     "*int32",
	"uint32":    "*uint32",
	"int64":     "*int64",
	"uint64":    "*uint64",
	"int":       "*int",
	"uint":      "*uint",
	"float32":   "*float32",
	"float64":   "*float64",
	"string":    "*string",
	"time.Time": "*time.Time",
}

var nullToSQLNull = map[string]string{
	"bool":      "sql.NullBool",
	"int8":      "*int8",
	"uint8":     "sql.NullByte",
	"int16":     "sql.NullInt16",
	"uint16":    "*uint16",
	"int32":     "sql.NullInt32",
	"uint32":    "*uint32",
	"int64":     "sql.NullInt64",
	"uint64":    "*uint64",
	"int":       "*int",
	"uint":      "*uint",
	"float32":   "sql.NullFloat64",
	"float64":   "sql.NullFloat64",
	"string":    "sql.NullString",
	"time.Time": "sql.NullTime",
}

// getFieldDataType get go data type name
func getFieldDataType(dataType string, isNull, disableNull, isNullToPointer, enableInt, enableIntegerInt bool) string {
	if enableInt {
		switch dataType {
		case "uint8", "uint16", "uint32":
			dataType = "uint"
		case "int8", "int16", "int32":
			dataType = "int"
		}
	}
	if enableIntegerInt {
		switch dataType {
		case "uint32":
			dataType = "uint"
		case "int32":
			dataType = "int"
		}
	}
	if !disableNull && isNull {
		cv := nullToSQLNull
		if isNullToPointer {
			cv = nullToPointer
		}
		if v, ok := cv[dataType]; ok {
			return v
		}
	}
	return dataType
}

var rJSONTag = regexp.MustCompile(`^.*?\[@.*?(?i:jsontag+):\s*(.*)\].*?`)
var rAffixJSONTag = regexp.MustCompile(`^.*?\[@.*?(affix+).*?\].*?`)

func jsonTag(comment string) string {
	match := rJSONTag.FindStringSubmatch(comment)
	if len(match) == 2 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func affixJSONTag(comment string) bool {
	match := rAffixJSONTag.FindStringSubmatch(comment)
	return len(match) == 2 && strings.TrimSpace(match[1]) == "affix"
}
