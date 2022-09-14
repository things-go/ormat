package ast

import (
	"fmt"
	"sort"
	"strings"
)

// Field of a struct
type Field struct {
	FieldName      string              // field name
	FieldType      string              // field type
	FieldComment   string              // field comment
	FieldTags      map[string][]string // field tags
	ColumnDataType string              // field column go type.
	ColumnName     string              // field column name in database
}

// AddFieldTag Add a tag
func (e *Field) AddFieldTag(k string, v string) *Field {
	if e.FieldTags == nil {
		e.FieldTags = make(map[string][]string)
	}
	e.FieldTags[k] = append(e.FieldTags[k], v)
	return e
}

func (e *Field) RemoveFieldTag(k, v string) *Field {
	if e.FieldTags != nil {
		tagsValues := e.FieldTags[k]
		for i, vv := range tagsValues {
			if vv == v {
				e.FieldTags[k] = append(tagsValues[:i], tagsValues[i+1:]...)
			}
		}
	}
	return e
}

// BuildLine build a field line
func (e *Field) BuildLine() string {
	var buf strings.Builder

	// field name
	buf.WriteString(e.FieldName)
	buf.WriteString(delimTab)

	// field type
	buf.WriteString(e.FieldType)
	buf.WriteString(delimTab)

	// field tags
	if len(e.FieldTags) > 0 {
		ks := make([]string, 0, len(e.FieldTags))
		for k := range e.FieldTags {
			ks = append(ks, k)
		}
		sort.Strings(ks)

		buf.WriteString("`")
		for i, v := range ks {
			buf.WriteString(fmt.Sprintf(`%v:"%v"`, v, strings.Join(e.FieldTags[v], ";")))
			if i != len(ks)-1 {
				buf.WriteString(" ")
			}
		}
		buf.WriteString("`")
	}

	comment := strings.ReplaceAll(e.FieldComment, "\n", ",")
	if comment != "" {
		buf.WriteString("// " + e.FieldComment)
	}
	return buf.String()
}
