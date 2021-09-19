package ast

import (
	"text/template"
)

// interval
const _interval = "\t"

const tableNameTpl = `
// TableName implement schema.Tabler interface
func (*{{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
`

var TableNameTpl, _ = template.New("tableNameTpl").Parse(tableNameTpl)
