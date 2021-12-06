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

// sqliteTable info
// sql: SELECT name FROM sqlite_master WHERE type='table' AND name !='sqlite_sequence'
type sqliteTable struct {
	Type     string `gorm:"column:type"`     // table type, etc. table, index
	Name     string `gorm:"column:name"`     // table name
	TblName  string `gorm:"column:tbl_name"` // belong table name
	RootPage int    `gorm:"column:rootpage"` // table
	SQL      string `gorm:"column:sql"`      // create table or index sql
}

//
// type sqliteKeys struct {
// 	NonUnique  int    `gorm:"column:Non_unique"`
// 	KeyName    string `gorm:"column:Key_name"`
// 	ColumnName string `gorm:"column:Column_name"`
// }

// sqliteColumn show full columns
type sqliteColumn struct {
	Cid       int     `gorm:"column:cid"`        // column ordinal position
	Name      string  `gorm:"column:name"`       // column name
	Type      string  `gorm:"column:type"`       // column type
	Pk        int     `gorm:"column:pk"`         // column pk
	NotNull   bool    `gorm:"column:notnull"`    // column not null
	DfltValue *string `gorm:"column:dflt_value"` // column default value
}

// sqliteForeignKey ...
// select table_schema,table_name,column_name,referenced_table_schema,referenced_table_name,referenced_column_name from INFORMATION_SCHEMA.KEY_COLUMN_USAGE
// where table_schema ='matrix' AND REFERENCED_TABLE_NAME IS NOT NULL AND TABLE_NAME = 'credit_card' ;
// foreignKey Foreign key of db info
type sqliteForeignKey struct {
	TableSchema           string `gorm:"column:table_schema"`            // Database of columns.列所在的数据库
	TableName             string `gorm:"column:table_name"`              // Data table of column.列所在的数据表
	ColumnName            string `gorm:"column:column_name"`             // column names.列名
	ReferencedTableSchema string `gorm:"column:referenced_table_schema"` // The database where the index is located.该索引所在的数据库
	ReferencedTableName   string `gorm:"column:referenced_table_name"`   // Affected tables . 该索引受影响的表
	ReferencedColumnName  string `gorm:"column:referenced_column_name"`  // Which column of the affected table.该索引受影响的表的哪一列
}

type SQLite struct{}

// GetDbInfo get database info
// 获取数据库信息
func (sf *SQLite) GetDatabase(db *gorm.DB, dbName string, tbNames ...string) (*view.Database, error) {
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
	// sort tables
	sort.Sort(view.Tables(tbInfos))
	return &view.Database{
		Name:   dbName,
		Tables: tbInfos,
	}, nil
}

// GetTables get all table name and comments
// 获取所有表及注释
func (sf *SQLite) GetTables(db *gorm.DB, dbName string, tbNames ...string) ([]view.TableAttribute, error) {
	var rows []sqliteTable

	err := db.Raw(`SELECT name FROM sqlite_master WHERE type='table' AND name !='sqlite_sequence'`).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]view.TableAttribute, 0, len(rows))
	for _, v := range rows {
		result = append(result, view.TableAttribute{Name: v.Name, Comment: "", CreateTableSQL: v.SQL})
	}
	return result, err
}

// GetTableColumns get table's column info.
// 获取表的所有列的信息
func (sf *SQLite) GetTableColumns(db *gorm.DB, dbName string, tb view.TableAttribute) (*view.Table, error) {
	var columnInfos []view.Column
	var columns []sqliteColumn
	var foreignKeys []sqliteForeignKey

	err := db.Raw("PRAGMA table_info(" + tb.Name + ")").Find(&columns).Error
	if err != nil {
		return nil, err
	}

	for _, v := range columns {
		columnInfo := view.Column{
			Name:            v.Name,
			OrdinalPosition: v.Cid,
			DataType:        getSqliteGoDataType(v.Name, v.Type),
			ColumnType:      v.Type, // TODO: ??
			IsNullable:      !v.NotNull,
			IsAutoIncrement: false,
			Default:         v.DfltValue,
			Comment:         "",
		}
		// TODO: 索引
		// if v.Pk == 1 {
		// 	columnInfo.Index = append(columnInfo.Index, view.Index{
		// 		KeyType: view.ColumnKeyTypePrimary,
		// 		IsMulti: false,
		// 	})
		// }

		// foreignKey
		// TODO: 外键
		columnInfo.ForeignKeys = fixSqliteForeignKey(foreignKeys, columnInfo.Name)

		columnInfos = append(columnInfos, columnInfo)
	}

	sort.Sort(view.Columns(columnInfos))
	return &view.Table{
		TableAttribute: tb,
		Columns:        columnInfos,
	}, nil
}

// GetCreateTableSQL get create table sql
func (*SQLite) GetCreateTableSQL(db *gorm.DB, tbName string) (string, error) {
	var row sqliteTable

	err := db.Raw("SELECT tbl_name, sql FROM sqlite_master WHERE type='table' AND name=?", tbName).
		Take(&row).Error
	return row.SQL, err
}

// fixForeignKey fix foreign key
// TODO: not implement
func fixSqliteForeignKey(vs []sqliteForeignKey, columnName string) []view.ForeignKey {
	var result []view.ForeignKey
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

// sqliteTypeDict Accurate matching type
var sqliteTypeDict = map[string]string{
	"integer": "int",
	"text":    "string",
	"real":    "float64",
	"numeric": "string",
	"blob":    "[]byte",
}

func getSqliteGoDataType(name, dataType string) string {
	dataType = getSqliteDataType(dataType)
	// filter special type
	switch name {
	case "created_at", "updated_at":
		if dataType == "string" {
			return "time.Time"
		}
	case "deleted_at":
		switch dataType {
		case "int8", "uint8", "int16", "uint16",
			"int32", "uint32", "int64", "uint64",
			"int", "uint":
			return "int64"
		}
	}
	return dataType

}

func getSqliteDataType(dataType string) string {
	dataType = strings.ToLower(dataType)
	selfDefineTypeMqlDicMap := config.GetTypeDefine()
	if v, ok := selfDefineTypeMqlDicMap[dataType]; ok {
		return v
	}
	if v, ok := sqliteTypeDict[dataType]; ok {
		return v
	}

	for _, v := range mysqlTypeDictMatchList {
		ok, _ := regexp.MatchString(v.Key, dataType)
		if ok {
			return v.Value
		}
	}
	panic(fmt.Sprintf("type (%v) not match in any way, need to add on (https://github.com/thinkgos/ormat/blob/master/view/model.go)", dataType))
	return ""
}
