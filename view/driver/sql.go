package driver

import (
	"errors"
	"fmt"

	"github.com/xwb1989/sqlparser"

	"github.com/things-go/ormat/view"
)

type SQL struct {
	CreateTableSQL   string
	CustomDefineType map[string]string
	table            *view.Table
}

func (sf *SQL) Parse() error {
	stmt, err := sqlparser.Parse(sf.CreateTableSQL)
	if err != nil {
		return err
	}
	switch stmt := stmt.(type) {
	case *sqlparser.DDL:
		tb := &view.Table{
			TableAttribute: view.TableAttribute{
				Name:    stmt.NewName.Name.String(),
				Comment: "",
				// CreateTableSQL: sf.CreateTableSQL,
			},
		}

		fmt.Printf("%#v\n", stmt.TableSpec.Columns[0])
		for _, column := range stmt.TableSpec.Columns {
			columnType := column.Type
			tb.Columns = append(tb.Columns, view.Column{
				Name:            column.Name.String(),
				OrdinalPosition: 0,
				DataType:        "",
				ColumnType:      "",
				IsNullable:      false,
				IsAutoIncrement: bool(columnType.Autoincrement),
				Default:         nil,
				Comment:         "",
				Index:           nil,
				ForeignKeys:     nil,
			})
		}

		// fmt.Printf("\n Indexes: \n")
		// for _, idx := range stmt.TableSpec.Indexes {
		// 	fmt.Printf("%#v\n", idx)
		// }
	default:
		return errors.New("sql is not DDL")
	}
	return nil

}

func (sf *SQL) GetDatabase() (*view.Database, error) {
	return nil, nil
}
func (sf *SQL) GetTables() ([]view.TableAttribute, error) {
	return nil, nil
}
func (sf *SQL) GetTableColumns(tb view.TableAttribute) (*view.Table, error) {
	return nil, nil
}
func (sf *SQL) GetCreateTableSQL(tbName string) (string, error) {
	return "", nil
}
