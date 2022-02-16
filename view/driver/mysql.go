package driver

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"gorm.io/gorm"

	"github.com/thinkgos/ormat/config"
	"github.com/thinkgos/ormat/view"
)

const Primary = "PRIMARY"

// mysqlTable mysql table info
// sql: SELECT * FROM information_schema.TABLES WHERE TABLE_SCHEMA={db_name}
type mysqlTable struct {
	Name    string `gorm:"column:TABLE_NAME"`    // table name, 表名
	Comment string `gorm:"column:TABLE_COMMENT"` // table comment, 表注释
}

// mysqlColumn mysql column info
// sql: SELECT * FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA`={dbName} AND `TABLE_NAME`={tbName}
type mysqlColumn struct {
	ColumnName             string  `gorm:"column:COLUMN_NAME"`      // column name
	OrdinalPosition        int     `gorm:"column:ORDINAL_POSITION"` // column ordinal position
	ColumnDefault          *string `gorm:"column:COLUMN_DEFAULT"`   // column default value.null mean not set.
	IsNullable             string  `gorm:"column:IS_NULLABLE"`      // column null or not, YEW/NO
	DataType               string  `gorm:"column:DATA_TYPE"`        // column data type(varchar)
	CharacterMaximumLength int64   `gorm:"column:CHARACTER_MAXIMUM_LENGTH"`
	CharacterOctetLength   int64   `gorm:"column:CHARACTER_OCTET_LENGTH"`
	NumericPrecision       int64   `gorm:"column:NUMERIC_PRECISION"`
	NumericScale           int64   `gorm:"column:NUMERIC_SCALE"`
	ColumnType             string  `gorm:"column:COLUMN_TYPE"`    // column type(varchar(64))
	ColumnKey              string  `gorm:"column:COLUMN_KEY"`     // column key, PRI/MUL
	Extra                  string  `gorm:"column:EXTRA"`          // extra (auto_increment)
	ColumnComment          string  `gorm:"column:COLUMN_COMMENT"` // column comment
}

// key index info
// sql: SHOW KEYS FROM {table_name}
type mysqlKey struct {
	Table      string `gorm:"column:Table"`        // 表名
	NonUnique  bool   `gorm:"column:Non_unique"`   // 不是唯一索引
	KeyName    string `gorm:"column:Key_name"`     // 索引关键字
	SeqInIndex int    `gorm:"column:Seq_in_index"` // 索引排序
	ColumnName string `gorm:"column:Column_name"`  // 索引列名
	IndexType  string `gorm:"column:Index_type"`   // 索引类型, BTREE
}

// mysqlForeignKey Foreign key of db table info . 表的外键信息
// sql: SELECT table_schema, table_name, column_name, referenced_table_schema, referenced_table_name, referenced_column_name
//		FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
//		WHERE table_schema={db_name} AND REFERENCED_TABLE_NAME IS NOT NULL AND TABLE_NAME={table_name}
type mysqlForeignKey struct {
	TableSchema           string `gorm:"column:table_schema"`            // Database of column.
	TableName             string `gorm:"column:table_name"`              // Data table of column.
	ColumnName            string `gorm:"column:column_name"`             // column names.
	ReferencedTableSchema string `gorm:"column:referenced_table_schema"` // The database where the index is located.
	ReferencedTableName   string `gorm:"column:referenced_table_name"`   // Affected tables .
	ReferencedColumnName  string `gorm:"column:referenced_column_name"`  // Which column of the affected table.
}

// mysqlCreateTable mysql show create table info
// sql: SHOW CREATE TABLE {tableName}
type mysqlCreateTable struct {
	Table string `gorm:"column:Table"`
	SQL   string `gorm:"column:Create Table"`
}

type MySQL struct{}

// GetDatabase get database information
func (sf *MySQL) GetDatabase(db *gorm.DB, dbName string, tbNames ...string) (*view.Database, error) {
	tables, err := sf.GetTables(db, dbName, tbNames...)
	if err != nil {
		return nil, err
	}

	tbInfos := make([]view.Table, 0, len(tables))
	for _, v := range tables {
		tbInfo, err := sf.GetTableColumns(db, dbName, v)
		if err != nil {
			return nil, err
		}
		tbInfos = append(tbInfos, *tbInfo)
	}
	sort.Sort(view.Tables(tbInfos))
	return &view.Database{
		Name:   dbName,
		Tables: tbInfos,
	}, nil
}

// GetTables get all table name and comments
func (sf *MySQL) GetTables(db *gorm.DB, dbName string, tbNames ...string) ([]view.TableAttribute, error) {
	var rows []mysqlTable

	err := db.Table("information_schema.TABLES").
		Scopes(func(db *gorm.DB) *gorm.DB {
			if len(tbNames) > 0 {
				db = db.Where("TABLE_NAME in (?)", tbNames)
			}
			return db.Where("TABLE_SCHEMA=?", dbName)
		}).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make([]view.TableAttribute, 0, len(rows))
	for _, v := range rows {
		createTableSQL, err := sf.GetCreateTableSQL(db, v.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, view.TableAttribute{Name: v.Name, Comment: v.Comment, CreateTableSQL: createTableSQL})
	}
	return result, err
}

// GetTableColumns get table's column info.
func (sf *MySQL) GetTableColumns(db *gorm.DB, dbName string, tb view.TableAttribute) (*view.Table, error) {
	var columnInfos []view.Column
	var columns []mysqlColumn
	var keys []mysqlKey
	var foreignKeys []mysqlForeignKey

	// get table column list
	err := db.Raw("SELECT * FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA`=? AND `TABLE_NAME`=?", dbName, tb.Name).
		Find(&columns).Error
	if err != nil {
		return nil, err
	}

	// get index key list
	err = db.Raw("SHOW KEYS FROM `" + tb.Name + "`").Find(&keys).Error
	if err != nil {
		return nil, err
	}

	keyNameCount := make(map[string]int)            // key name count
	ColumnNameMapKey := make(map[string][]mysqlKey) // column name map key
	for _, v := range keys {
		keyNameCount[v.KeyName]++
		ColumnNameMapKey[v.ColumnName] = append(ColumnNameMapKey[v.ColumnName], v)
	}

	// get table column foreign keys
	err = db.Raw(`
			SELECT table_schema, table_name, column_name, referenced_table_schema, referenced_table_name, referenced_column_name
			FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE 
			WHERE table_schema=? AND REFERENCED_TABLE_NAME IS NOT NULL AND TABLE_NAME=?`, dbName, tb.Name).
		Find(&foreignKeys).Error
	if err != nil {
		return nil, err
	}

	for _, v := range columns {
		ci := view.Column{
			Name:            v.ColumnName,
			OrdinalPosition: v.OrdinalPosition,
			DataType:        getMysqlGoDataType(v.ColumnType),
			ColumnType:      v.ColumnType,
			IsNullable:      strings.EqualFold(v.IsNullable, "YES"),
			IsAutoIncrement: v.Extra == "auto_increment",
			Default:         v.ColumnDefault,
			Comment:         v.ColumnComment,
		}

		// column keys
		if columnKeys, ok := ColumnNameMapKey[v.ColumnName]; ok {
			for _, vv := range columnKeys {
				// non unique, normal index
				kk := view.ColumnKeyTypeNormalIndex
				// primary or unique
				if !vv.NonUnique {
					if strings.EqualFold(vv.KeyName, Primary) { // primary key
						kk = view.ColumnKeyTypePrimary
					} else {
						kk = view.ColumnKeyTypeUniqueKey // unique index
					}
				}
				ci.Index = append(ci.Index, view.Index{
					KeyType:    kk,
					IsMulti:    keyNameCount[vv.KeyName] > 1,
					KeyName:    vv.KeyName,
					SeqInIndex: vv.SeqInIndex,
					IndexType:  vv.IndexType,
				})
			}
		}

		// foreignKey
		ci.ForeignKeys = fixForeignKey(foreignKeys, ci.Name)

		columnInfos = append(columnInfos, ci)
	}

	sort.Sort(view.Columns(columnInfos))
	return &view.Table{
		TableAttribute: tb,
		Columns:        columnInfos,
	}, nil
}

// GetCreateTableSQL get create table sql
func (sf *MySQL) GetCreateTableSQL(db *gorm.DB, tbName string) (string, error) {
	var ct mysqlCreateTable

	err := db.Raw("SHOW CREATE TABLE `" + tbName + "`").
		Take(&ct).Error
	if err != nil {
		return "", err
	}
	return ct.SQL, err
}

// fixForeignKey fix foreign key
// TODO: not implement
func fixForeignKey(vs []mysqlForeignKey, columnName string) []view.ForeignKey {
	result := make([]view.ForeignKey, 0, len(vs))
	for _, v := range vs {
		if strings.EqualFold(v.ColumnName, columnName) {
			result = append(result, view.ForeignKey{
				TableName:  v.ReferencedTableName,
				ColumnName: v.ReferencedColumnName,
			})
		}
	}
	return result
}

func getMysqlGoDataType(columnType string) string {
	selfDefineTypeMqlDicMap := config.GetTypeDefine()
	if v, ok := selfDefineTypeMqlDicMap[columnType]; ok {
		return v
	}
	for _, v := range typeDictMatchList {
		ok, _ := regexp.MatchString(v.Key, columnType)
		if ok {
			return v.Value
		}
	}
	panic(fmt.Sprintf("type (%v) not match in any way, need to add on (https://github.com/thinkgos/ormat/blob/master/view/model.go)", columnType))
	return ""
}

type dictMatchKv struct {
	Key   string
	Value string
}

// \b([(]\d+[)])? 匹配0个或1个(\d+)
var typeDictMatchList = []dictMatchKv{
	{`^(tinyint)\b[(]1[)] unsigned`, "bool"},
	{`^(tinyint)\b[(]1[)]`, "bool"},
	{`^(tinyint)\b[(]\d+[)] unsigned`, "uint8"},
	{`^(tinyint)\b[(]\d+[)]`, "int8"},
	{`^(smallint)\b[(]\d+[)] unsigned`, "uint16"},
	{`^(smallint)\b[(]\d+[)]`, "int16"},
	{`^(mediumint)\b[(]\d+[)] unsigned`, "uint32"},
	{`^(mediumint)\b[(]\d+[)]`, "int32"},
	{`^(int)\b[(]\d+[)] unsigned`, "uint32"},
	{`^(int)\b[(]\d+[)]`, "int32"},
	{`^(integer)\b[(]\d+[)] unsigned`, "uint32"},
	{`^(integer)\b[(]\d+[)]`, "int32"},
	{`^(bigint)\b[(]\d+[)] unsigned`, "uint64"},
	{`^(bigint)\b[(]\d+[)]`, "int64"},
	{`^(float)\b[(]\d+,\d+[)] unsigned`, "float32"},
	{`^(float)\b[(]\d+,\d+[)]`, "float32"},
	{`^(double)\b([(]\d+,\d+[)])? unsigned`, "float64"},
	{`^(double)\b([(]\d+,\d+[)])?`, "float64"},
	{`^(char)\b[(]\d+[)]`, "string"},
	{`^(varchar)\b[(]\d+[)]`, "string"},
	{`^(datetime)\b([(]\d+[)])?`, "time.Time"},
	{`^(date)\b([(]\d+[)])?`, "datatypes.Date"},
	{`^(timestamp)\b([(]\d+[)])?`, "time.Time"},
	{`^(time)\b([(]\d+[)])?`, "time.Time"},
	{`^(text)\b([(]\d+[)])?`, "string"},
	{`^(tinytext)\b([(]\d+[)])?`, "string"},
	{`^(mediumtext)\b([(]\d+[)])?`, "string"},
	{`^(longtext)\b([(]\d+[)])?`, "string"},
	{`^(blob)\b([(]\d+[)])?`, "[]byte"},
	{`^(tinyblob)\b([(]\d+[)])?`, "[]byte"},
	{`^(mediumblob)\b([(]\d+[)])?`, "[]byte"},
	{`^(longblob)\b([(]\d+[)])?`, "[]byte"},
	{`^(bit)\b[(]\d+[)]`, "[]uint8"},
	{`^(json)\b`, "datatypes.JSON"},
	{`^(enum)\b[(](.)+[)]`, "string"},
	{`^(decimal)\b[(]\d+,\d+[)]`, "string"},
	{`^(binary)\b[(]\d+[)]`, "[]byte"},
	{`^(varbinary)\b[(]\d+[)]`, "[]byte"},
}
