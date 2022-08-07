package view

import (
	"bytes"
	"strings"

	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/things-go/ormat/utils"
	"github.com/things-go/ormat/view/ast"
)

// DBModel Implement the interface to acquire database information.
type DBModel interface {
	GetDatabase(db *gorm.DB, dbName string, tbNames ...string) (*Database, error)
	GetTables(db *gorm.DB, dbName string, tbNames ...string) ([]TableAttribute, error)
	GetTableColumns(db *gorm.DB, dbName string, tb TableAttribute) (*Table, error)
	GetCreateTableSQL(db *gorm.DB, tbName string) (string, error)
}

type WebTag struct {
	Kind    string `yaml:"kind" json:"kind"`
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
	db      *gorm.DB
	dbName  string
	tbNames []string
}

// New view instance
func New(db *gorm.DB, m DBModel, c Config, dbName string, tbNames ...string) *View {
	return &View{c, m, db, dbName, tbNames}
}

// GetDbFile ast file
func (sf *View) GetDbFile(pkgName string) ([]ast.File, error) {
	dbInfo, err := sf.GetDatabase(sf.db, sf.dbName, sf.tbNames...)
	if err != nil {
		return nil, err
	}

	files := make([]ast.File, 0, len(dbInfo.Tables))
	for _, sct := range sf.GetTableStruct(dbInfo.Tables) {
		file := new(ast.File).
			SetName(sct.GetTableName() + ".go").
			SetPackageName(pkgName).
			SetOutColumnName(sf.IsOutColumnName).
			AddStruct(sct)

		files = append(files, *file)
	}
	return files, nil
}

// GetDBCreateTableSQLContent get all table's create table sql content
func (sf *View) GetDBCreateTableSQLContent() ([]byte, error) {
	tbSqls, err := sf.GetTables(sf.db, sf.dbName, sf.tbNames...)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	for _, vv := range tbSqls {
		buf.WriteString("# " + vv.Name + " " + strings.ReplaceAll(vv.Comment, "\n", "\n# ") + "\n" + vv.CreateTableSQL + ";\n\n")
	}
	return buf.Bytes(), nil
}

// GetTableStruct get table struct
func (sf *View) GetTableStruct(tables []Table) []ast.Struct {
	scts := make([]ast.Struct, 0, len(tables))
	for _, tb := range tables {
		sct := new(ast.Struct).
			SetName(utils.CamelCase(tb.Name, sf.EnableLint)).
			SetComment(tb.Comment).
			AddFields(sf.getColumnFields(tables, tb.Columns)...).
			SetTableName(tb.Name).
			SetCreatTableSQL(tb.CreateTableSQL).
			EnableOutSQL(sf.IsOutSQL)

		scts = append(scts, *sct)
	}
	return scts
}

// getColumnFields Get table column's field
func (sf *View) getColumnFields(tables []Table, cols []Column) []ast.Field {
	fields := make([]ast.Field, 0, len(cols))
	for _, v := range cols {
		var field ast.Field

		fieldName := utils.CamelCase(v.Name, sf.EnableLint)
		fieldType := getFieldDataType(v.DataType, v.IsNullable, sf.DisableNull, sf.IsNullToPoint, sf.EnableInt, sf.EnableIntegerInt, sf.EnableBoolInt)
		if fieldName == "DeletedAt" &&
			(v.DataType == "int64" ||
				v.DataType == "uint64" ||
				v.DataType == "uint32" ||
				v.DataType == "int32" ||
				v.DataType == "uint16" ||
				v.DataType == "int16" ||
				v.DataType == "uint8" ||
				v.DataType == "int8" ||
				v.DataType == "uint" ||
				v.DataType == "int") {
			fieldType = "soft_delete.DeletedAt"
		}
		field.SetName(fieldName).
			SetType(fieldType).
			SetComment(v.Comment).
			SetColumnName(v.Name)

		sf.fixFieldTags(&field, v)

		fields = append(fields, field)

		if sf.IsForeignKey && len(v.ForeignKeys) > 0 {
			fks := sf.getForeignKeyField(tables, v)
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
				field.SetName(name + "List").
					SetType("[]" + name)
			} else {
				field.SetName(name).
					SetType(name)
			}
			field.SetComment(comment).
				AddTag(tagDb, "joinForeignKey:"+col.Name).
				AddTag(tagDb, "foreignKey:"+v.ColumnName)

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
	if tagDb != "" {
		// not simple output
		field.AddTag(tagDb, "column:"+ci.Name)
		columnType := "type:" + ci.ColumnType
		field.AddTag(tagDb, columnType)

		if ci.IsAutoIncrement {
			field.AddTag(tagDb, "autoIncrement:true")
		}
		if !ci.IsNullable {
			field.AddTag(tagDb, "not null")
		}
		// default tag
		if ci.Default != nil {
			dflt := "default:''"
			if *ci.Default != "" {
				dflt = "default:" + *ci.Default
			}
			field.AddTag(tagDb, dflt)
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
					field.RemoveTag(tagDb, columnType)
				}
				if v1.IsMulti {
					if vv == "primaryKey" {
						vv += ";"
					} else {
						vv += ","
					}
					vv += "priority:" + cast.ToString(v1.SeqInIndex)
				}
				field.AddTag(tagDb, vv)
			}
		}
		if sf.IsCommentTag && field.GetComment() != "" {
			comment := strings.TrimSpace(field.GetComment())
			comment = strings.ReplaceAll(comment, ";", ",")
			comment = strings.ReplaceAll(comment, "`", "'")
			comment = strings.ReplaceAll(comment, `"`, `\"`)
			comment = strings.ReplaceAll(comment, "\n", " ")
			comment = strings.ReplaceAll(comment, "\r\n", " ")
			field.AddTag(tagDb, "comment:"+comment)
		}
	}

	// web tag
	fixFieldWebTags(field, ci.Name, sf.WebTags, sf.EnableLint)
}

func fixFieldWebTags(field *ast.Field, name string, webTags []WebTag, enableLint bool) {
	for _, v := range webTags {
		vv := ""
		if v.Tag == "json" {
			if vv = jsonTag(field.Comment); vv != "" {
				field.AddTag(v.Tag, vv)
				return
			}
		}

		switch v.Kind {
		case "smallCamelCase":
			vv = utils.SmallCamelCase(name, enableLint)
		case "camelCase":
			vv = utils.CamelCase(name, enableLint)
		case "snakeCase":
			vv = utils.SnakeCase(name, enableLint)
		case "kebab":
			vv = utils.Kebab(name, enableLint)
		}

		if vv != "" {
			if v.HasOmit {
				vv += ",omitempty"
			}
			if v.Tag == "json" && affixJSONTag(field.Comment) {
				vv += ",string"
			}
			field.AddTag(v.Tag, vv)
		}
	}
}
