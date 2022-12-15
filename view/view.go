package view

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/things-go/ormat/pkg/consts"
	"github.com/things-go/ormat/pkg/matcher"
	"github.com/things-go/ormat/pkg/utils"
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
	GetDatabase() (*Database, error)
	GetTableAttributes() ([]TableAttribute, error)
	GetTables(tb TableAttribute) (*Table, error)
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
	IsForeignKey     bool     `yaml:"isForeignKey" json:"isForeignKey"`         // 输出外键
	IsCommentTag     bool     `yaml:"isCommentTag" json:"isCommentTag"`         // 注释同时放入tag标签中
}

// View information
type View struct {
	Config
	DBModel
}

// New view instance
func New(m DBModel, c Config) *View {
	return &View{c, m}
}

// GetDbFile ast file
func (sf *View) GetDbFile(pkgName string) ([]*ast.File, error) {
	dbInfo, err := sf.GetDatabase()
	if err != nil {
		return nil, err
	}

	files := make([]*ast.File, 0, len(dbInfo.Tables))
	for _, tb := range dbInfo.Tables {
		structName := utils.CamelCase(tb.Name, sf.EnableLint)
		structComment := ast.IntoComment(tb.Comment, "...", "\n", "\r\n// ")
		structFields := sf.intoColumnFields(dbInfo.Tables, tb.Columns)
		tableName := tb.Name
		abbrTableName := ast.IntoAbbrTableName(tableName)
		protoMessageFields, protoEnum := ast.ParseProtobuf(structName, tableName, structFields)
		structs := []*ast.Struct{
			{
				StructName:         structName,
				StructComment:      structComment,
				StructFields:       structFields,
				TableName:          tableName,
				AbbrTableName:      abbrTableName,
				CreateTableSQL:     tb.CreateTableSQL,
				ProtoMessageFields: protoMessageFields,
				ProtoEnum:          protoEnum,
			},
		}
		files = append(files, &ast.File{
			Version:     consts.Version,
			Filename:    tb.Name,
			PackageName: pkgName,
			Imports:     ast.IntoImports(structs),
			Structs:     structs,
		})
	}
	return files, nil
}

func (sf *View) GetDBCreateTableSQL() (*ast.SqlFile, error) {
	tbAttributes, err := sf.GetTableAttributes()
	if err != nil {
		return nil, err
	}
	tbAttrs := make([]ast.TableAttribute, 0, len(tbAttributes))
	for _, v := range tbAttributes {
		tbAttrs = append(tbAttrs, ast.TableAttribute{
			Name:           v.Name,
			Comment:        strings.ReplaceAll(v.Comment, "\n", "\n-- "),
			CreateTableSQL: v.CreateTableSQL,
		})
	}
	return &ast.SqlFile{
		Version: consts.Version,
		Tables:  tbAttrs,
	}, nil
}

// GetDBCreateTableSQLContent get all table's create table sql content
func (sf *View) GetDBCreateTableSQLContent() ([]byte, error) {
	tbAttributes, err := sf.GetTableAttributes()
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	for _, ta := range tbAttributes {
		buf.WriteString(
			"-- " + ta.Name + " " + strings.ReplaceAll(ta.Comment, "\n", "\n-- ") + "\n" +
				ta.CreateTableSQL + ";\n\n",
		)
	}
	return buf.Bytes(), nil
}

// intoColumnFields Get table column's field
func (sf *View) intoColumnFields(tables []*Table, cols []*Column) []ast.Field {
	fields := make([]ast.Field, 0, len(cols))
	for _, col := range cols {
		fieldName := utils.CamelCase(col.Name, sf.EnableLint)
		fieldType := getFieldDataType(col.ColumnGoType, col.IsNullable, sf.DisableNull, sf.IsNullToPoint, sf.EnableInt, sf.EnableIntegerInt, sf.EnableBoolInt)
		if fieldName == "DeletedAt" &&
			(col.ColumnGoType == "int64" ||
				col.ColumnGoType == "uint64" ||
				col.ColumnGoType == "uint32" ||
				col.ColumnGoType == "int32" ||
				col.ColumnGoType == "uint16" ||
				col.ColumnGoType == "int16" ||
				col.ColumnGoType == "uint8" ||
				col.ColumnGoType == "int8" ||
				col.ColumnGoType == "uint" ||
				col.ColumnGoType == "int") {
			fieldType = "soft_delete.DeletedAt"
		}

		field := ast.Field{
			FieldName:    fieldName,
			FieldType:    fieldType,
			FieldComment: ast.IntoComment(col.Comment, "", "\n", ","),
			FieldTag:     "",
			ColumnGoType: col.ColumnGoType,
			ColumnName:   col.Name,
		}
		fieldTags := ast.NewFieldTags()
		sf.fixFieldTags(fieldTags, &field, col)
		field.FieldTag = fieldTags.IntoFieldTag()
		fields = append(fields, field)

		if sf.IsForeignKey && len(col.ForeignKeys) > 0 {
			fks := sf.intoForeignKeyField(tables, col)
			fields = append(fields, fks...)
		}
	}
	return fields
}

// intoForeignKeyField Get information about foreign key of table column field
// TODO: not implement.
func (sf *View) intoForeignKeyField(tables []*Table, col *Column) (fks []ast.Field) {
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
			fieldTags := ast.NewFieldTags().
				AddTagValue(tagDb, "joinForeignKey:"+col.Name).
				AddTagValue(tagDb, "foreignKey:"+v.ColumnName)

			fixFieldWebTags(fieldTags, &field, v.TableName, sf.WebTags, sf.EnableLint)
			field.FieldTag = fieldTags.IntoFieldTag()
			fks = append(fks, field)
		}
	}
	return
}

func (*View) getColumnsKeyMulti(tables []*Table, tableName, col string) (isMulti bool, isFind bool, notes string) {
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

func (sf *View) fixFieldTags(fieldTags *ast.FieldTags, field *ast.Field, ci *Column) {
	tagDb := sf.DbTag
	if tagDb == "" {
		tagDb = "gorm"
	}

	columnType := "type:" + ci.ColumnType
	// 输出db标签
	filedTagValues := ast.NewFiledTagValues().
		SetSeparate(";").
		AddValue("column:" + ci.Name).
		AddValue(columnType)

	if ci.IsAutoIncrement {
		filedTagValues.AddValue("autoIncrement:true")
	}
	if !ci.IsNullable {
		filedTagValues.AddValue("not null")
	}
	// default tag
	if ci.Default != nil {
		dflt := "default:''"
		if *ci.Default != "" {
			dflt = "default:" + *ci.Default
		}
		filedTagValues.AddValue(dflt)
	}

	for _, index := range ci.Index {
		var vv string

		switch index.KeyType {
		// case ColumnsDefaultKey:
		case ColumnKeyTypePrimary:
			vv = "primaryKey"
		case ColumnKeyTypeUniqueKey:
			vv = "uniqueIndex:" + index.KeyName
		case ColumnKeyTypeNormalIndex:
			vv = "index:" + index.KeyName
			if index.KeyName == "sort" { // 兼容 gorm 本身 sort 标签
				vv = "index"
			}
			if index.IndexType == "FULLTEXT" {
				vv += ",class:FULLTEXT"
			}
		}
		if vv != "" {
			// NOTE: 主要是整型主键,gorm在自动迁移时没有在mysql上加上auto_increment
			if vv == "primaryKey" && ci.IsAutoIncrement {
				filedTagValues.RemoveValue(columnType)
			}
			if index.IsMulti {
				if vv == "primaryKey" {
					vv += ";"
				} else {
					vv += ","
				}
				vv += "priority:" + strconv.FormatInt(int64(index.SeqInIndex), 10)
			}
			filedTagValues.AddValue(vv)
		}
	}
	if sf.IsCommentTag && field.FieldComment != "" {
		comment := strings.TrimSpace(field.FieldComment)
		comment = strings.ReplaceAll(comment, ";", ",")
		comment = strings.ReplaceAll(comment, "`", "'")
		comment = strings.ReplaceAll(comment, `"`, `\"`)
		comment = strings.ReplaceAll(comment, "\r\n", " ")
		comment = strings.ReplaceAll(comment, "\n", " ")
		filedTagValues.AddValue("comment:" + comment)
	}
	fieldTags.Add(tagDb, filedTagValues)

	// web tag
	fixFieldWebTags(fieldTags, field, ci.Name, sf.WebTags, sf.EnableLint)
}

func fixFieldWebTags(fieldTags *ast.FieldTags, field *ast.Field, columnName string, webTags []WebTag, enableLint bool) {
	intoWebTagName := func(kind, columnName string, enableLint bool) string {
		vv := ""
		switch kind {
		case WebTagSmallCamelCase:
			vv = utils.SmallCamelCase(columnName, enableLint)
		case WebTagCamelCase:
			vv = utils.CamelCase(columnName, enableLint)
		case WebTagSnakeCase:
			vv = utils.SnakeCase(columnName, enableLint)
		case WebTagKebab:
			vv = utils.Kebab(columnName, enableLint)
		}
		return vv
	}

	for _, v := range webTags {
		if v.Tag == "json" {
			if vv := matcher.JsonTag(field.FieldComment); vv != "" {
				fieldTags.Add(
					v.Tag,
					ast.NewFiledTagValues().
						SetSeparate(",").
						AddValue(vv),
				)
				continue
			}
		}
		vv := intoWebTagName(v.Kind, columnName, enableLint)
		if vv == "" {
			continue
		}
		filedTagValue := ast.NewFiledTagValues().
			SetSeparate(",").
			AddValue(vv)
		if v.HasOmit {
			filedTagValue.AddValue("omitempty")
		}
		if v.Tag == "json" && matcher.HasAffixJSONTag(field.FieldComment) {
			filedTagValue.AddValue("string")
		}
		fieldTags.Add(v.Tag, filedTagValue)
	}
}
