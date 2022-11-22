package view

import (
	"bytes"
	"strings"

	"github.com/spf13/cast"

	"github.com/things-go/ormat/utils"
	"github.com/things-go/ormat/view/ast"
)

const (
	WebTagSmallCamelCase = "smallCamelCase"
	WebTagCamelCase      = "camelCase"
	WebTagSnakeCase      = "snakeCase"
	WebTagKebab          = "kebab"
)

// DBModel Implement the interface to acquire database information.
type DBModel interface {
	GetDatabase(dbName string, tbNames ...string) (*Database, error)
	GetTables(dbName string, tbNames ...string) ([]TableAttribute, error)
	GetTableColumns(dbName string, tb TableAttribute) (*Table, error)
	GetCreateTableSQL(tbName string) (string, error)
}

type WebTag struct {
	Kind    string `yaml:"kind" json:"kind"` // support smallCamelCase, camelCase, snakeCase, kebab
	Tag     string `yaml:"tag" json:"tag"`
	HasOmit bool   `yaml:"hasOmit" json:"hasOmit"`
}

type Config struct {
	DbTag            string   `yaml:"dbTag" json:"dbTag"`                       // db标签, 默认gorm
	WebTags          []WebTag `yaml:"webTags" json:"webTags"`                   // web tags 标签列表
	EnableLint       bool     `yaml:"enableLint" json:"enableLint"`             // 使能lint, id -> ID
	DisableNull      bool     `yaml:"disableNull" json:"disableNull"`           // 不输出字段为null指针或sql.Nullxxx类型
	EnableInt        bool     `yaml:"enableInt" json:"enableInt"`               // 使能int8,uint8,int16,uint16,int32,uint32输出为int, uint
	EnableIntegerInt bool     `yaml:"enableIntegerInt" json:"enableIntegerInt"` // 使能int32,uint32输出为int, uint
	EnableBoolInt    bool     `yaml:"enableBoolInt" json:"enableBoolInt"`       // 使能bool输出int
	IsNullToPoint    bool     `yaml:"isNullToPoint" json:"isNullToPoint"`       // 是否字段为null时输出指针类型
	IsOutSQL         bool     `yaml:"isOutSQL" json:"isOutSQL"`                 // 是否输出创建表的SQL
	IsOutColumnName  bool     `yaml:"isOutColumnName" json:"isOutColumnName"`   // 是否输出表的列名, 默认不输出
	IsForeignKey     bool     `yaml:"isForeignKey" json:"isForeignKey"`         // 输出外键
	IsCommentTag     bool     `yaml:"isCommentTag" json:"isCommentTag"`         // 注释同时放入tag标签中
}

// View information
type View struct {
	Config
	DBModel
	dbName  string
	tbNames []string
}

// New view instance
func New(m DBModel, c Config, dbName string, tbNames ...string) *View {
	return &View{c, m, dbName, tbNames}
}

// GetDbFile ast file
func (sf *View) GetDbFile(pkgName string) ([]*ast.File, error) {
	dbInfo, err := sf.GetDatabase(sf.dbName, sf.tbNames...)
	if err != nil {
		return nil, err
	}

	files := make([]*ast.File, 0, len(dbInfo.Tables))
	for _, sct := range sf.GetTableStruct(dbInfo.Tables) {
		files = append(files, &ast.File{
			Filename:        sct.TableName,
			PackageName:     pkgName,
			Imports:         make(map[string]string),
			Structs:         []ast.Struct{*sct},
			IsOutColumnName: sf.IsOutColumnName,
		})
	}
	return files, nil
}

// GetDBCreateTableSQLContent get all table's create table sql content
func (sf *View) GetDBCreateTableSQLContent() ([]byte, error) {
	tbSqls, err := sf.GetTables(sf.dbName, sf.tbNames...)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	for _, vv := range tbSqls {
		buf.WriteString("-- " + vv.Name + " " + strings.ReplaceAll(vv.Comment, "\n", "\n-- ") + "\n" + vv.CreateTableSQL + ";\n\n")
	}
	return buf.Bytes(), nil
}

// GetTableStruct get table struct
func (sf *View) GetTableStruct(tables []Table) []*ast.Struct {
	structs := make([]*ast.Struct, 0, len(tables))
	for _, tb := range tables {
		structs = append(structs, &ast.Struct{
			StructName:     utils.CamelCase(tb.Name, sf.EnableLint),
			StructComment:  tb.Comment,
			StructFields:   sf.getColumnFields(tables, tb.Columns),
			TableName:      tb.Name,
			CreateTableSQL: tb.CreateTableSQL,
		})
	}
	return structs
}

// getColumnFields Get table column's field
func (sf *View) getColumnFields(tables []Table, cols []Column) []ast.Field {
	fields := make([]ast.Field, 0, len(cols))
	for _, col := range cols {
		fieldName := utils.CamelCase(col.Name, sf.EnableLint)
		fieldType := getFieldDataType(col.DataType, col.IsNullable, sf.DisableNull, sf.IsNullToPoint, sf.EnableInt, sf.EnableIntegerInt, sf.EnableBoolInt)
		if fieldName == "DeletedAt" &&
			(col.DataType == "int64" ||
				col.DataType == "uint64" ||
				col.DataType == "uint32" ||
				col.DataType == "int32" ||
				col.DataType == "uint16" ||
				col.DataType == "int16" ||
				col.DataType == "uint8" ||
				col.DataType == "int8" ||
				col.DataType == "uint" ||
				col.DataType == "int") {
			fieldType = "soft_delete.DeletedAt"
		}

		field := ast.Field{
			FieldName:      fieldName,
			FieldType:      fieldType,
			FieldComment:   col.Comment,
			FieldTags:      make(map[string]*ast.FieldTagValue),
			ColumnDataType: col.DataType,
			ColumnName:     col.Name,
		}
		sf.fixFieldTags(&field, col)

		fields = append(fields, field)

		if sf.IsForeignKey && len(col.ForeignKeys) > 0 {
			fks := sf.getForeignKeyField(tables, col)
			fields = append(fields, fks...)
		}
	}
	return fields
}

// getForeignKeyField Get information about foreign key of table column field
// TODO: not implement.
func (sf *View) getForeignKeyField(tables []Table, col Column) (fks []ast.Field) {
	tagDb := sf.DbTag
	if tagDb == "" {
		tagDb = "gorm"
	}

	for _, v := range col.ForeignKeys {
		isMulti, found, comment := sf.getColumnsKeyMulti(tables, v.TableName, v.ColumnName)
		if found {
			var field ast.Field

			name := utils.CamelCase(v.TableName, sf.EnableLint)
			if isMulti {
				field.FieldName = name + "List"
				field.FieldType = "[]" + name
			} else {
				field.FieldName = name
				field.FieldType = name
			}
			field.FieldComment = comment
			field.AddFieldTagValue(tagDb, "joinForeignKey:"+col.Name).
				AddFieldTagValue(tagDb, "foreignKey:"+v.ColumnName)

			fixFieldWebTags(&field, v.TableName, sf.WebTags, sf.EnableLint)
			fks = append(fks, field)
		}
	}
	return
}

func (*View) getColumnsKeyMulti(tables []Table, tableName, col string) (isMulti bool, isFind bool, notes string) {
	for _, v := range tables {
		if strings.EqualFold(v.Name, tableName) {
			for _, v1 := range v.Columns {
				if strings.EqualFold(v1.Name, col) {
					for _, v2 := range v1.Index {
						switch v2.KeyType {
						case ColumnKeyTypePrimary, ColumnKeyTypeUniqueKey:
							if !v2.IsMulti { // 唯一索引
								return false, true, v.Comment
							}
						case ColumnKeyTypeNormalIndex: // index key. 复合索引
							isMulti = true
						}
					}
					return true, true, v.Comment
				}
			}
			break
		}
	}
	return false, false, ""
}

func (sf *View) fixFieldTags(field *ast.Field, ci Column) {
	tagDb := sf.DbTag
	if tagDb == "" {
		tagDb = "gorm"
	}

	// 输出db标签
	// not simple output
	columnType := "type:" + ci.ColumnType
	filedTagValue := ast.NewFiledTagValue().
		SetSeparate(";").
		AddValue("column:" + ci.Name).
		AddValue(columnType)

	if ci.IsAutoIncrement {
		filedTagValue.AddValue("autoIncrement:true")
	}
	if !ci.IsNullable {
		filedTagValue.AddValue("not null")
	}
	// default tag
	if ci.Default != nil {
		dflt := "default:''"
		if *ci.Default != "" {
			dflt = "default:" + *ci.Default
		}
		filedTagValue.AddValue(dflt)
	}

	for _, v1 := range ci.Index {
		var vv string

		switch v1.KeyType {
		// case ColumnsDefaultKey:
		case ColumnKeyTypePrimary:
			vv = "primaryKey"
		case ColumnKeyTypeUniqueKey:
			vv = "uniqueIndex:" + v1.KeyName
		case ColumnKeyTypeNormalIndex:
			vv = "index:" + v1.KeyName
			// 兼容 gorm 本身 sort 标签
			if v1.KeyName == "sort" {
				vv = "index"
			}
			if v1.IndexType == "FULLTEXT" {
				vv += ",class:FULLTEXT"
			}
		}
		if vv != "" {
			// NOTE: 主要是整型主键,gorm在自动迁移时没有在mysql上加上auto_increment
			if vv == "primaryKey" && ci.IsAutoIncrement {
				filedTagValue.RemoveValue(columnType)
			}
			if v1.IsMulti {
				if vv == "primaryKey" {
					vv += ";"
				} else {
					vv += ","
				}
				vv += "priority:" + cast.ToString(v1.SeqInIndex)
			}
			filedTagValue.AddValue(vv)
		}
	}
	if sf.IsCommentTag && field.FieldComment != "" {
		comment := strings.TrimSpace(field.FieldComment)
		comment = strings.ReplaceAll(comment, ";", ",")
		comment = strings.ReplaceAll(comment, "`", "'")
		comment = strings.ReplaceAll(comment, `"`, `\"`)
		comment = strings.ReplaceAll(comment, "\n", " ")
		comment = strings.ReplaceAll(comment, "\r\n", " ")
		filedTagValue.AddValue("comment:" + comment)
	}
	field.AddFieldTag(tagDb, filedTagValue)

	// web tag
	fixFieldWebTags(field, ci.Name, sf.WebTags, sf.EnableLint)
}

func fixFieldWebTags(field *ast.Field, name string, webTags []WebTag, enableLint bool) {
	for _, v := range webTags {
		filedTagValue := ast.NewFiledTagValue().
			SetSeparate(",")
		vv := ""
		if v.Tag == "json" {
			if vv = jsonTag(field.FieldComment); vv != "" {
				filedTagValue.AddValue(vv)
				return
			}
		}

		switch v.Kind {
		case WebTagSmallCamelCase:
			vv = utils.SmallCamelCase(name, enableLint)
		case WebTagCamelCase:
			vv = utils.CamelCase(name, enableLint)
		case WebTagSnakeCase:
			vv = utils.SnakeCase(name, enableLint)
		case WebTagKebab:
			vv = utils.Kebab(name, enableLint)
		}

		if vv != "" {
			filedTagValue.AddValue(vv)
			if v.HasOmit {
				filedTagValue.AddValue("omitempty")
			}
			if v.Tag == "json" && affixJSONTag(field.FieldComment) {
				filedTagValue.AddValue("string")
			}
			field.AddFieldTag(v.Tag, filedTagValue)
		}
	}
}
