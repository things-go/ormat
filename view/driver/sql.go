package driver

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"github.com/xwb1989/sqlparser"

	"github.com/things-go/ormat/view"
)

type SQL struct {
	CreateTableSQL   string
	CustomDefineType map[string]string
	table            *view.Table
}

func (sf *SQL) hasParse() bool {
	return sf.table != nil
}

func (sf *SQL) Parse() error {
	if sf.hasParse() {
		return nil
	}
	statement, err := sqlparser.Parse(sf.CreateTableSQL)
	if err != nil {
		return err
	}
	switch stmt := statement.(type) {
	case *sqlparser.DDL:
		if stmt.Action != sqlparser.CreateStr {
			return errors.New("不是创建表语句")
		}
		if stmt.TableSpec == nil {
			return errors.New("未解析到任何字段")
		}

		tb := &view.Table{
			TableAttribute: view.TableAttribute{
				Name:           stmt.NewName.Name.String(),
				Comment:        "",
				CreateTableSQL: sf.CreateTableSQL,
			},
		}
		// ENGINE=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci comment='我是注释'
		tbOptions := strings.Split(stmt.TableSpec.Options, " ")
		for _, option := range tbOptions {
			keyValue := strings.Split(option, "=")
			if len(keyValue) >= 2 {
				switch keyValue[0] {
				case "ENGINE":
				case "charset":
				case "collate":
				case "comment":
					tb.Comment = strings.ReplaceAll(keyValue[1], "'", "")
				}
			}
		}

		columnFiledMapping := make(map[string]*view.Column)
		for i := 0; i < len(stmt.TableSpec.Columns); i++ {
			var defaultValue *string
			var comment string

			column := stmt.TableSpec.Columns[i]
			columnType := column.Type

			if columnType.Default != nil {
				val := string(columnType.Default.Val)
				defaultValue = &val
			}
			if columnType.Comment != nil {
				comment = string(columnType.Comment.Val)
			}

			ct, dt, err := intoDataTypeAndColumnType(&columnType)
			if err != nil {
				return err
			}
			ci := &view.Column{
				Name:            column.Name.String(),
				OrdinalPosition: i + 1,
				DataType:        dt,
				ColumnType:      ct,
				IsNullable:      !bool(columnType.NotNull),
				IsAutoIncrement: bool(columnType.Autoincrement),
				Default:         defaultValue,
				Comment:         comment,
				Index:           nil,
				ForeignKeys:     nil,
			}
			tb.Columns = append(tb.Columns, ci)
			columnFiledMapping[column.Name.String()] = ci
		}
		for _, indexes := range stmt.TableSpec.Indexes {
			keyType := view.ColumnKeyNone
			switch indexes.Info.Type {
			case "primary key":
				keyType = view.ColumnKeyTypePrimary
			case "unique key":
				keyType = view.ColumnKeyTypeUniqueKey
			case "unique":
				keyType = view.ColumnKeyTypeUnique
			case "key":
				keyType = view.ColumnKeyTypeNormalIndex
			}
			indexType := "BTREE"
			for _, option := range indexes.Options {
				if option.Name == "using" {
					indexType = option.Using
				}
			}
			isMulti := len(indexes.Columns) > 1
			for i := 0; i < len(indexes.Columns); i++ {
				col := indexes.Columns[i]
				ci := columnFiledMapping[col.Column.String()]
				if ci == nil {
					break
				}
				seqInIndex := 0
				if isMulti {
					seqInIndex = i + 1
				}
				ci.Index = append(ci.Index, view.Index{
					KeyType:    keyType,
					KeyName:    indexes.Info.Name.String(),
					IsMulti:    isMulti,
					SeqInIndex: seqInIndex,
					IndexType:  indexType,
				})
			}
		}
		sf.table = tb
	default:
		return errors.New("不是DDL语句")
	}
	return nil
}

func (sf *SQL) GetDatabase() (*view.Database, error) {
	err := sf.Parse()
	if err != nil {
		return nil, err
	}
	return &view.Database{
		Name:   "",
		Tables: []*view.Table{sf.table},
	}, nil
}
func (sf *SQL) GetTables() ([]view.TableAttribute, error) {
	err := sf.Parse()
	if err != nil {
		return nil, err
	}
	return []view.TableAttribute{sf.table.TableAttribute}, nil
}
func (sf *SQL) GetTableColumns(tb view.TableAttribute) (*view.Table, error) {
	err := sf.Parse()
	if err != nil {
		return nil, err
	}
	return sf.table, nil
}
func (sf *SQL) GetCreateTableSQL(tbName string) (string, error) {
	return sf.CreateTableSQL, nil
}

func intoDataTypeAndColumnType(columnType *sqlparser.ColumnType) (string, string, error) {
	var ct, dt string

	toInt := func(l *sqlparser.SQLVal) int {
		length := 0
		if l != nil {
			length = cast.ToInt(string(l.Val))
		}
		return length
	}

	isUnsigned := bool(columnType.Unsigned)
	switch columnType.Type {
	case "tinyint", "smallint", "mediumint", "int", "integer", "bigint":
		// {`^(tinyint)\b[(]1[)] unsigned`, "bool"},
		// {`^(tinyint)\b[(]1[)]`, "bool"},
		// {`^(tinyint)\b([(]\d+[)])? unsigned`, "uint8"},
		// {`^(tinyint)\b([(]\d+[)])?`, "int8"},
		// {`^(smallint)\b([(]\d+[)])? unsigned`, "uint16"},
		// {`^(smallint)\b([(]\d+[)])?`, "int16"},
		// {`^(mediumint)\b([(]\d+[)])? unsigned`, "uint32"},
		// {`^(mediumint)\b([(]\d+[)])?`, "int32"},
		// {`^(int)\b([(]\d+[)])? unsigned`, "uint32"},
		// {`^(int)\b([(]\d+[)])?`, "int32"},
		// {`^(integer)\b([(]\d+[)])? unsigned`, "uint32"},
		// {`^(integer)\b([(]\d+[)])?`, "int32"},
		// {`^(bigint)\b([(]\d+[)])? unsigned`, "uint64"},
		// {`^(bigint)\b([(]\d+[)])?`, "int64"},

		dataTypeMapping := map[string]struct {
			unsigned string
			signed   string
		}{
			"tinyint":   {"uint8", "int8"},
			"smallint":  {"uint16", "int16"},
			"mediumint": {"uint32", "int32"},
			"int":       {"uint32", "int32"},
			"integer":   {"uint32", "int32"},
			"bigint":    {"uint64", "int64"},
		}

		toColumnType := func(targetColumnType string, length int, isUnsigned bool) string {
			if length > 0 {
				if isUnsigned {
					return fmt.Sprintf("%s(%d) unsigned", targetColumnType, length)
				} else {
					return fmt.Sprintf("%s(%d)", targetColumnType, length)

				}
			} else {
				if isUnsigned {
					return fmt.Sprintf("%s unsigned", targetColumnType)
				} else {
					return targetColumnType
				}
			}
		}

		length := toInt(columnType.Length)
		if columnType.Type == "tinyint" && length == 1 {
			dt = "bool"
		} else {
			vv := dataTypeMapping[columnType.Type]
			if isUnsigned {
				dt = vv.unsigned
			} else {
				dt = vv.signed
			}
		}
		ct = toColumnType(columnType.Type, length, isUnsigned)
	case "float", "double", "decimal":
		// {`^(float)\b([(]\d+,\d+[)])? unsigned`, "float32"},
		// {`^(float)\b([(]\d+,\d+[)])?`, "float32"},
		// {`^(double)\b([(]\d+,\d+[)])? unsigned`, "float64"},
		// {`^(double)\b([(]\d+,\d+[)])?`, "float64"},
		// {`^(decimal)\b[(]\d+,\d+[)]`, "string"},
		dataTypeMapping := map[string]string{
			"float":   "float32",
			"double":  "float64",
			"decimal": "string",
		}
		toColumnType := func(targetColumnType string, length, scale int) string {
			if length > 0 {
				return fmt.Sprintf("%s(%d,%d)", targetColumnType, length, scale)
			} else {
				return targetColumnType
			}
		}
		length, scale := toInt(columnType.Length), toInt(columnType.Scale)
		dt = dataTypeMapping[columnType.Type]
		ct = toColumnType(columnType.Type, length, scale)

	case "char", "varchar",
		"text", "tinytext", "mediumtext", "longtext",
		"date", "datetime", "timestamp", "time",
		"blob", "tinyblob", "mediumblob", "longblob",
		"binary", "varbinary",
		"bit":
		// {`^(char)\b[(]\d+[)]`, "string"},
		// {`^(varchar)\b[(]\d+[)]`, "string"},
		// {`^(text)\b([(]\d+[)])?`, "string"},
		// {`^(tinytext)\b([(]\d+[)])?`, "string"},
		// {`^(mediumtext)\b([(]\d+[)])?`, "string"},
		// {`^(longtext)\b([(]\d+[)])?`, "string"},
		// {`^(date)\b([(]\d+[)])?`, "datatypes.Date"},
		// {`^(datetime)\b([(]\d+[)])?`, "time.Time"},
		// {`^(timestamp)\b([(]\d+[)])?`, "time.Time"},
		// {`^(time)\b([(]\d+[)])?`, "time.Time"},
		// {`^(blob)\b([(]\d+[)])?`, "[]byte"},
		// {`^(tinyblob)\b([(]\d+[)])?`, "[]byte"},
		// {`^(mediumblob)\b([(]\d+[)])?`, "[]byte"},
		// {`^(longblob)\b([(]\d+[)])?`, "[]byte"},
		// {`^(binary)\b[(]\d+[)]`, "[]byte"},
		// {`^(varbinary)\b[(]\d+[)]`, "[]byte"},
		// {`^(bit)\b[(]\d+[)]`, "[]uint8"},
		dataTypeMapping := map[string]string{
			"char":       "string",
			"varchar":    "string",
			"text":       "string",
			"tinytext":   "string",
			"mediumtext": "string",
			"longtext":   "string",
			"date":       "datatypes.Date",
			"datetime":   "time.Time",
			"timestamp":  "time.Time",
			"time":       "time.Time",
			"blob":       "[]byte",
			"tinyblob":   "[]byte",
			"mediumblob": "[]byte",
			"longblob":   "[]byte",
			"binary":     "[]byte",
			"varbinary":  "[]byte",
			"bit":        "[]uint8",
		}
		toColumnType := func(targetColumnType string, length int) string {
			if length > 0 {
				return fmt.Sprintf("%s(%d)", targetColumnType, length)
			} else {
				return targetColumnType
			}
		}
		length := toInt(columnType.Length)
		dt = dataTypeMapping[columnType.Type]
		ct = toColumnType(columnType.Type, length)
	case "enum":
		// {`^(enum)\b[(](.)+[)]`, "string"},
		dt = "string"
		ct = "enum(" + strings.Join(columnType.EnumValues, ",") + ")"
	case "json":
		// {`^(json)\b`, "datatypes.JSON"},
		dt = "datatypes.JSON"
		ct = columnType.Type
	default:
		return ct, dt, fmt.Errorf("not support column type(%s)", columnType.Type)
	}

	return ct, dt, nil
}
