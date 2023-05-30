package view

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/exp/slices"
)

// ColumnKeyType column key type
type ColumnKeyType int

const (
	ColumnKeyNone            ColumnKeyType = iota // default
	ColumnKeyTypePrimary                          // primary key
	ColumnKeyTypeNormalIndex                      // normal index key
	ColumnKeyTypeUniqueKey                        // unique key
	ColumnKeyTypeUnique                           // unique
)

// Database database information
type Database struct {
	Name   string   // database name, 数据库名
	Tables []*Table // table information, 表信息
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
	Columns []*Column // column information
	Indexes []*Index  // index information
}

// IntoColumnDefinedSQL 转换为定义的字段 SQL
func (c *Table) IntoIndexDefinedSQL() []string {
	m := make(map[string][]*Index)
	for _, index := range c.Indexes {
		m[index.KeyName] = append(m[index.KeyName], index)
	}

	pk := make([]string, 0, len(m))
	ks := make([]string, 0, len(m))
	for k, indexes := range m {
		sort.Sort(IndexSlice(indexes))
		b := strings.Builder{}
		b.Grow(64)

		index0 := indexes[0]
		switch index0.KeyType {
		case ColumnKeyTypePrimary:
			b.WriteString("PRIMARY KEY")
		case ColumnKeyTypeNormalIndex:
			b.WriteString("KEY")
		case ColumnKeyTypeUniqueKey:
			b.WriteString("UNIQUE KEY")
		case ColumnKeyTypeUnique:
			b.WriteString("UNIQUE")
		default:
			continue
		}
		if index0.KeyType != ColumnKeyTypePrimary {
			b.WriteString(" ")
			b.WriteString("`")
			b.WriteString(k)
			b.WriteString("`")
		}
		b.WriteString(" ")
		b.WriteString("(")
		first := true
		for _, v := range indexes {
			if !first {
				b.WriteString(", ")
			}
			first = false
			b.WriteString("`")
			b.WriteString(v.ColumnName)
			b.WriteString("`")
		}
		b.WriteString(")")
		if index0.IndexType != "" {
			b.WriteString(" USING ")
			b.WriteString(index0.IndexType)
		}
		if index0.KeyType == ColumnKeyTypePrimary {
			pk = append(pk, b.String())
		} else {
			ks = append(ks, b.String())
		}
	}
	sort.Strings(pk)
	sort.Strings(ks)
	return append(pk, ks...)
}

// Index database index/unique_index list
type Index struct {
	KeyType     ColumnKeyType // key type, 索引类型
	KeyName     string        // index key name, 索引名称
	IsComposite bool          // composite key, 是否为复合索引
	SeqInIndex  int           // union index sequence in index, 复合索引中的序列
	ColumnName  string        // column name, 列名
	IndexType   string        // index type, 索引类型(比如: BTREE, HASH, FULLTEXT)
}

// ForeignKey Foreign key
type ForeignKey struct {
	TableName  string // Affected tables .
	ColumnName string // Which column of the affected table
}

// Column column information
type Column struct {
	Name            string       // column name
	OrdinalPosition int          // column ordinal position
	ColumnGoType    string       // column data go type(string,int...)
	ColumnType      string       // column type(varchar(256)...)
	IsNullable      bool         // column is null or not
	IsAutoIncrement bool         // column auto increment or not
	Default         *string      // default value
	Comment         string       // column comment
	Index           []*Index     // index list
	ForeignKeys     []ForeignKey // Foreign key list
}

// IntoColumnDefinedSQL 转换为定义的字段 SQL
func (c *Column) IntoColumnDefinedSQL() string {
	b := strings.Builder{}
	b.Grow(64)

	b.WriteString(c.ColumnType)
	if !c.IsNullable {
		b.WriteString(" ")
		b.WriteString("NOT NULL")
	}
	if c.IsAutoIncrement {
		b.WriteString(" ")
		b.WriteString("AUTO_INCREMENT")
	} else {
		dv := ""
		if c.IsNullable {
			dv = "DEFAULT NULL"
			if c.Default != nil && *c.Default != "null" {
				dv = fmt.Sprintf("DEFAULT '%s'", *c.Default)
			}
		} else {
			if c.Default != nil {
				dv = fmt.Sprintf("DEFAULT '%s'", *c.Default)
			} else if slices.Contains(
				[]string{
					"bool",
					"int8", "uint8", "int16", "uint16",
					"int32", "uint32", "int64", "uint64",
					"int", "uint", "float32", "float64",
				},
				c.ColumnGoType) {
				dv = "DEFAULT '0'"
			} else if c.ColumnGoType == "string" {
				dv = "DEFAULT ''"
			}
		}
		if dv != "" {
			b.WriteString(" ")
			b.WriteString(dv)
		}
	}

	return b.String()
}

type TableSlice []*Table

func (t TableSlice) Len() int           { return len(t) }
func (t TableSlice) Less(i, j int) bool { return t[i].Name < t[j].Name }
func (t TableSlice) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

type IndexSlice []*Index

func (t IndexSlice) Len() int           { return len(t) }
func (t IndexSlice) Less(i, j int) bool { return t[i].SeqInIndex < t[j].SeqInIndex }
func (t IndexSlice) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

type ColumnSlice []*Column

func (t ColumnSlice) Len() int           { return len(t) }
func (t ColumnSlice) Less(i, j int) bool { return t[i].OrdinalPosition < t[j].OrdinalPosition }
func (t ColumnSlice) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

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

// intoFieldDataType get go data type name
func intoFieldDataType(columnGoType string, isNullable, disableNullToPointer, enableInt, enableIntegerInt, enableBoolInt bool) string {
	if enableInt {
		switch columnGoType {
		case "uint8", "uint16", "uint32":
			columnGoType = "uint"
		case "int8", "int16", "int32":
			columnGoType = "int"
		}
	}
	if enableIntegerInt {
		switch columnGoType {
		case "uint32":
			columnGoType = "uint"
		case "int32":
			columnGoType = "int"
		}
	}
	if enableBoolInt && columnGoType == "bool" {
		columnGoType = "int"
	}

	if isNullable {
		cv := nullToPointer
		if disableNullToPointer {
			cv = nullToSQLNull
		}
		if v, ok := cv[columnGoType]; ok {
			return v
		}
	}
	return columnGoType
}

var goTypeToAssistType = map[string]string{
	"bool":           "Bool",
	"int8":           "Int8",
	"uint8":          "Uint8",
	"int16":          "Int16",
	"uint16":         "Uint16",
	"int32":          "Int32",
	"uint32":         "Uint32",
	"int64":          "Int64",
	"uint64":         "Uint64",
	"int":            "Int",
	"uint":           "Uint",
	"float32":        "Float32",
	"float64":        "Float64",
	"decimal":        "Decimal",
	"string":         "String",
	"[]byte":         "Byte",
	"datatypes.Date": "Time",
	"time.Time":      "Time",
}

var d = regexp.MustCompile(`^(decimal)\b[(]\d+,\d+[)]`)

func IsDecimal(t string) bool {
	return d.MatchString(t)
}

// intoFieldAssistType get go data assist type name
func intoFieldAssistType(columnGoType, columnType string, enableInt, enableIntegerInt, enableBoolInt bool) string {
	if enableInt {
		switch columnGoType {
		case "uint8", "uint16", "uint32":
			columnGoType = "uint"
		case "int8", "int16", "int32":
			columnGoType = "int"
		}
	}
	if enableIntegerInt {
		switch columnGoType {
		case "uint32":
			columnGoType = "uint"
		case "int32":
			columnGoType = "int"
		}
	}
	if enableBoolInt && columnGoType == "bool" {
		columnGoType = "int"
	}
	if IsDecimal(columnType) {
		columnGoType = "decimal"
	}
	if t, ok := goTypeToAssistType[columnGoType]; ok {
		return t
	}
	return "Field"
}
