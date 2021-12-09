package ast

import (
	"strings"
)

// Struct define a struct
type Struct struct {
	Name           string  // struct name
	Comment        string  // struct comment
	Fields         []Field // struct field list
	TableName      string  // table name
	CreateTableSQL string  // create table SQL
	outSQL         bool
}

// SetName set the  name
func (s *Struct) SetName(name string) *Struct {
	s.Name = name
	return s
}

// GetName get the struct name
func (s *Struct) GetName() string { return s.Name }

// SetComment set the comment
func (s *Struct) SetComment(comment string) *Struct {
	s.Comment = comment
	return s
}

// SetComment set the comment
func (s *Struct) GetComment() string { return s.Comment }

// AddFields Add one or more fields
func (s *Struct) AddFields(e ...Field) *Struct {
	s.Fields = append(s.Fields, e...)
	return s
}

// SetTableName set the table name in database
func (s *Struct) SetTableName(name string) *Struct {
	s.TableName = name
	return s
}

// GetTableName get the table name in database
func (s *Struct) GetTableName() string { return s.TableName }

// SetCreatTableSQL set create table sql
func (s *Struct) SetCreatTableSQL(sql string) *Struct {
	s.CreateTableSQL = sql
	return s
}

// GetCreatTableSQL get create table sql
func (s *Struct) GetCreatTableSQL() string { return s.CreateTableSQL }

// EnableOutSQL enable out sql
func (s *Struct) EnableOutSQL(b bool) *Struct {
	s.outSQL = b
	return s
}

func (s *Struct) BuildTableNameTemplate() string {
	type Tpl struct {
		TableName  string
		StructName string
	}

	var buf strings.Builder

	_ = TableNameTpl.Execute(&buf, Tpl{
		TableName:  s.TableName,
		StructName: s.Name,
	})
	return buf.String()
}

// BuildLines Get the result data.获取结果数据
func (s *Struct) BuildLines() []string {
	var lines []string

	if s.outSQL {
		lines = append(lines,
			"/* sql",
			s.CreateTableSQL,
			"sql */",
			delimLF,
		)
	}

	comment := s.Comment
	if comment != "" {
		comment = strings.ReplaceAll(strings.TrimSpace(comment), "\n", "\r\n// ")
	} else {
		comment = "..."
	}
	comment = "// " + s.Name + " " + comment
	lines = append(lines,
		comment,
		"type\t"+s.Name+"\tstruct {",
	)

	// field every line
	mp := make(map[string]struct{}, len(s.Fields))
	for _, v := range s.Fields {
		name := v.GetName()
		if _, ok := mp[name]; !ok {
			mp[name] = struct{}{}
			lines = append(lines, "\t\t"+v.BuildLine())
		}
	}
	lines = append(lines, "}")
	return lines
}
