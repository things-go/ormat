package ast

import (
	"fmt"
	"sort"
	"strings"
)

// Field of a var
type Field struct {
	Name       string              // field name
	Type       string              // field type
	Comment    string              // field comment
	Tags       map[string][]string // field tags
	ColumnName string              // field column name in database
}

// SetName set field name
func (e *Field) SetName(name string) *Field {
	e.Name = name
	return e
}

// GetName get field name
func (e *Field) GetName() string {
	return e.Name
}

// SetType set field type
func (e *Field) SetType(tp string) *Field {
	e.Type = tp
	return e
}

// GetType get field type
func (e *Field) GetType() string { return e.Type }

// SetComment set field comment
func (e *Field) SetComment(comment string) *Field {
	e.Comment = comment
	return e
}

// GetComment get field comment
func (e *Field) GetComment() string { return e.Comment }

// AddTag Add a tag
func (e *Field) AddTag(k string, v string) *Field {
	if e.Tags == nil {
		e.Tags = make(map[string][]string)
	}
	e.Tags[k] = append(e.Tags[k], v)
	return e
}

// SetColumnName set field column name
func (e *Field) SetColumnName(name string) *Field {
	e.ColumnName = name
	return e
}

// GetColumnName get field column name
func (e *Field) GetColumnName() string {
	return e.ColumnName
}

// BuildLine build a field line
func (e *Field) BuildLine() string {
	var buf strings.Builder

	// field name
	buf.WriteString(e.Name)
	buf.WriteString(delimTab)

	// field type
	buf.WriteString(e.Type)
	buf.WriteString(delimTab)

	// field tags
	if len(e.Tags) > 0 {
		ks := make([]string, 0, len(e.Tags))
		for k := range e.Tags {
			ks = append(ks, k)
		}
		sort.Strings(ks)

		buf.WriteString("`")
		for i, v := range ks {
			buf.WriteString(fmt.Sprintf(`%v:"%v"`, v, strings.Join(e.Tags[v], ";")))
			if i != len(ks)-1 {
				buf.WriteString(" ")
			}
		}
		buf.WriteString("`")
	}

	comment := strings.ReplaceAll(e.Comment, "\n", ",")
	if comment != "" {
		buf.WriteString("// " + e.Comment)
	}
	return buf.String()
}
