package view

import (
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/things-go/ormat/pkg/consts"
	"github.com/things-go/ormat/pkg/matcher"
	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/view/ast"
)

const (
	TagSmallCamelCase = "smallCamelCase"
	TagCamelCase      = "camelCase"
	TagSnakeCase      = "snakeCase"
	TagKebab          = "kebab"
)

// DBModel Implement the interface to acquire database information.
type DBModel interface {
	GetDatabase() (*Database, error)
	GetTableAttributes() ([]TableAttribute, error)
	GetTables(tb TableAttribute) (*Table, error)
	GetCreateTableSQL(tbName string) (string, error)
}

type Config struct {
	DbTag              string            `yaml:"dbTag" json:"dbTag"`                         // db标签,默认gorm
	Tags               map[string]string `yaml:"tags" json:"tags"`                           // tags标签列表, support smallCamelCase, camelCase, snakeCase, kebab
	EnableInt          bool              `yaml:"enableInt" json:"enableInt"`                 // 使能int8,uint8,int16,uint16,int32,uint32输出为int,uint
	EnableIntegerInt   bool              `yaml:"enableIntegerInt" json:"enableIntegerInt"`   // 使能int32,uint32输出为int,uint
	EnableBoolInt      bool              `yaml:"enableBoolInt" json:"enableBoolInt"`         // 使能bool输出int
	DisableNullToPoint bool              `yaml:"isNullToPoint" json:"isNullToPoint"`         // 禁用字段为null时输出指针类型,将输出为sql.Nullxx
	DisableCommentTag  bool              `yaml:"disableCommentTag" json:"disableCommentTag"` // 禁用注释放入tag标签中
	EnableForeignKey   bool              `yaml:"enableForeignKey" json:"enableForeignKey"`   // 输出外键
	HasColumn          bool              `yaml:"hasColumn" json:"hasColumn"`                 // 是否输出字段
	SkipColumns        []string          `yaml:"skipColumns" json:"skipColumns"`             // 忽略输出字段, 格式 table.column
	Package            string            `yaml:"package" json:"package"`                     // 包名
	Options            map[string]string `yaml:"options" json:"options"`                     // 选项
	HasHelper          bool              `yaml:"hasHelper" json:"hasHelper"`                 // 是否输出 proto 帮助
	EnableGogo         bool              `yaml:"enableGogo" json:"enableGogo"`               // 使能用 gogo proto (仅 hasHelper = true 有效果)
	EnableSea          bool              `yaml:"enableSea" json:"enableSea"`                 // 使能用 seaql(仅 hasHelper = true 有效果)
}

func InitFlagSetForConfig(s *flag.FlagSet, cc *Config) {
	s.StringVarP(&cc.DbTag, "dbTag", "k", "gorm", "db标签")
	s.StringToStringVarP(&cc.Tags, "tags", "K", map[string]string{"json": TagSnakeCase}, "tags标签,类型支持[smallCamelCase,camelCase,snakeCase,kebab]")
	s.BoolVarP(&cc.EnableInt, "enableInt", "e", false, "使能int8,uint8,int16,uint16,int32,uint32输出为int,uint")
	s.BoolVarP(&cc.EnableIntegerInt, "enableIntegerInt", "E", false, "使能int32,uint32输出为int,uint")
	s.BoolVarP(&cc.EnableBoolInt, "enableBoolInt", "b", false, "使能bool输出int")
	s.BoolVarP(&cc.DisableNullToPoint, "disableNullToPoint", "B", false, "禁用字段为null时输出指针类型,将输出为sql.Nullxx")
	s.BoolVarP(&cc.DisableCommentTag, "disableCommentTag", "j", false, "禁用注释放入tag标签中")
	s.BoolVarP(&cc.EnableForeignKey, "enableForeignKey", "J", false, "使用外键")
	s.BoolVar(&cc.HasColumn, "hasColumn", false, "是否输出字段")
	s.StringSliceVar(&cc.SkipColumns, "skipColumns", nil, "忽略输出字段(仅 hasColumn = true 有效), 格式 table.column(只作用于指定表字段) 或  column(作用于所有表)")
	s.StringVar(&cc.Package, "package", "", "package name")
	s.StringToStringVar(&cc.Options, "options", nil, "options key value")

	s.BoolVar(&cc.HasHelper, "hasHelper", false, "是否输出 proto 帮助")
	s.BoolVar(&cc.EnableGogo, "enableGogo", false, "使能用 gogo proto (仅 hasHelper = true 有效)")
	s.BoolVar(&cc.EnableSea, "enableSea", false, "使能用 seaql (仅 hasHelper = true 有效)")
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
	skipColumns := make(map[string]struct{})
	for _, column := range sf.SkipColumns {
		skipColumns[utils.SnakeCase(column)] = struct{}{}
	}

	dbInfo, err := sf.GetDatabase()
	if err != nil {
		return nil, err
	}

	files := make([]*ast.File, 0, len(dbInfo.Tables))
	for _, tb := range dbInfo.Tables {
		structName := utils.CamelCase(tb.Name)
		structComment := ast.IntoComment(tb.Comment, "...", "\n", "\r\n// ")
		structFields := sf.intoColumnFields(tb.Name, dbInfo.Tables, tb.Columns, skipColumns)
		tableName := tb.Name
		abbrTableName := ast.IntoAbbrTableName(tableName)
		protoMessageFields := ast.ParseProtobuf(structFields, sf.EnableGogo, sf.EnableSea)
		structs := []*ast.Struct{
			{
				StructName:         structName,
				StructComment:      structComment,
				StructFields:       structFields,
				TableName:          tableName,
				AbbrTableName:      abbrTableName,
				CreateTableSQL:     tb.CreateTableSQL,
				ProtoMessageFields: protoMessageFields,
			},
		}
		files = append(files, &ast.File{
			Version:     consts.Version,
			Filename:    tb.Name,
			PackageName: pkgName,
			Imports:     ast.IntoImports(structs),
			Structs:     structs,
			Package:     sf.Package,
			Options:     sf.Options,
			HasColumn:   sf.HasColumn,
			HasHelper:   sf.HasHelper,
		})
	}
	return files, nil
}

func (sf *View) GetSqlFile() (*ast.SqlFile, error) {
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

// intoColumnFields Get table column's field
func (sf *View) intoColumnFields(tbName string, tables []*Table, cols []*Column, skipColumns map[string]struct{}) []ast.Field {
	fields := make([]ast.Field, 0, len(cols))
	for _, col := range cols {
		fieldName := utils.CamelCase(col.Name)
		fieldType := intoFieldDataType(col.ColumnGoType, col.IsNullable, sf.DisableNullToPoint, sf.EnableInt, sf.EnableIntegerInt, sf.EnableBoolInt)
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

		isSkipColumn := false
		if _, isSkipColumn = skipColumns[col.Name]; !isSkipColumn {
			_, isSkipColumn = skipColumns[tbName+"."+col.Name]
		}

		field := ast.Field{
			FieldName:    fieldName,
			FieldType:    fieldType,
			FieldComment: ast.IntoComment(col.Comment, "", "\n", ","),
			FieldTag:     "",
			IsNullable:   col.IsNullable,
			IsTimestamp:  col.ColumnGoType == "time.Time",
			ColumnGoType: col.ColumnGoType,
			ColumnName:   col.Name,
			Type:         col.IntoSqlDefined(),
			IsSkipColumn: isSkipColumn,
		}
		fieldTags := ast.NewFieldTags()
		sf.fixFieldTags(fieldTags, &field, col)
		field.FieldTag = fieldTags.IntoFieldTag()
		fields = append(fields, field)

		if sf.EnableForeignKey && len(col.ForeignKeys) > 0 {
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

			name := utils.CamelCase(v.TableName)
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

			fixFieldTags(fieldTags, &field, v.TableName, sf.Tags)
			field.FieldTag = fieldTags.IntoFieldTag()
			fks = append(fks, field)
		}
	}
	return
}

func (*View) getColumnsKeyMulti(tables []*Table, tableName, col string) (isMulti bool, isFind bool, notes string) {
	for _, tb := range tables {
		if strings.EqualFold(tb.Name, tableName) {
			for _, column := range tb.Columns {
				if strings.EqualFold(column.Name, col) {
					for _, idx := range column.Index {
						switch idx.KeyType {
						case ColumnKeyTypePrimary, ColumnKeyTypeUniqueKey:
							if !idx.IsMulti { // 唯一索引
								return false, true, tb.Comment
							}
						case ColumnKeyTypeNormalIndex: // index key. 复合索引
							isMulti = true
						}
					}
					return true, true, tb.Comment
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
	if !sf.DisableCommentTag && field.FieldComment != "" {
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
	fixFieldTags(fieldTags, field, ci.Name, sf.Tags)
}

func fixFieldTags(fieldTags *ast.FieldTags, field *ast.Field, columnName string, tags map[string]string) {
	intoWebTagName := func(kind, columnName string) string {
		vv := ""
		switch kind {
		case TagSmallCamelCase:
			vv = utils.SmallCamelCase(columnName)
		case TagCamelCase:
			vv = utils.CamelCase(columnName)
		case TagSnakeCase:
			vv = utils.SnakeCase(columnName)
		case TagKebab:
			vv = utils.Kebab(columnName)
		}
		return vv
	}

	for tag, kind := range tags {
		if tag == "json" {
			if vv := matcher.JsonTag(field.FieldComment); vv != "" {
				fieldTags.Add(
					tag,
					ast.NewFiledTagValues().
						SetSeparate(",").
						AddValue(vv),
				)
				continue
			}
		}
		vv := intoWebTagName(kind, columnName)
		if vv == "" {
			continue
		}
		filedTagValue := ast.NewFiledTagValues().
			SetSeparate(",").
			AddValue(vv).
			AddValue("omitempty")
		if tag == "json" && matcher.HasAffixJSONTag(field.FieldComment) {
			filedTagValue.AddValue("string")
		}
		fieldTags.Add(tag, filedTagValue)
	}
}
