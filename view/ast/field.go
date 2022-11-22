package ast

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

// FieldTagValue field tag value
type FieldTagValue struct {
	Sep   string   // tag separate like [",", ";", " "]
	Value []string // tag value
}

// NewFiledTagValue new FieldTagValue instance with default separate `,`.
func NewFiledTagValue() *FieldTagValue {
	return &FieldTagValue{
		Sep: ",",
	}
}

// SetSeparate set separate.
func (f *FieldTagValue) SetSeparate(sep string) *FieldTagValue {
	f.Sep = sep
	return f
}

// AddValue add value
func (f *FieldTagValue) AddValue(value string) *FieldTagValue {
	f.Value = append(f.Value, value)
	return f
}

// RemoveValue remove a tag
func (e *FieldTagValue) RemoveValue(value string) *FieldTagValue {
	tagsValues := e.Value
	for i, vv := range tagsValues {
		if vv == value {
			e.Value = append(tagsValues[:i], tagsValues[i+1:]...)
		}
	}
	return e
}

// IsEmpty value count is empty
func (f *FieldTagValue) IsEmpty() bool { return len(f.Value) == 0 }

// Field of a struct
//
//	  FieldName  FieldType  FieldTags                               FieldComment
//
//		 |        |            |                                         |
//	     v        v            v                                         v
//		Foo      int      `json:"foo,omitempty yaml:"foo,omitempty"` // 我是一个注释
//
// FieldTags
//
//	key   v1 sep v2 ...
//	 |     | |
//	 v     v v
//	json:"foo,omitempty"
type Field struct {
	FieldName      string                    // field name
	FieldType      string                    // field type
	FieldComment   string                    // field comment
	FieldTags      map[string]*FieldTagValue // field tags
	ColumnDataType string                    // field column go type.
	ColumnName     string                    // field column name in database
}

// AddFieldTag Add a tag
func (e *Field) GetFieldTagValue(key string) *FieldTagValue {
	if e.FieldTags == nil {
		e.FieldTags = make(map[string]*FieldTagValue)
	}
	fieldTagValue, ok := e.FieldTags[key]
	if !ok {
		fieldTagValue = NewFiledTagValue()
		e.FieldTags[key] = fieldTagValue
	}
	return fieldTagValue
}

// AddFieldTag add a tag value
func (e *Field) AddFieldTag(key string, tagValue *FieldTagValue) *Field {
	if tagValue != nil {
		e.FieldTags[key] = tagValue
	}
	return e
}

// RemoveFieldTag remove a tag value
func (e *Field) RemoveFieldTag(key string) *Field {
	delete(e.FieldTags, key)
	return e
}

// SetFieldTagSep set a tag separate
func (e *Field) SetFieldTagSep(key, sep string) *Field {
	fieldTagValue := e.GetFieldTagValue(key)
	fieldTagValue.SetSeparate(sep)
	return e
}

// AddFieldTagValue add a tag value
func (e *Field) AddFieldTagValue(key, value string) *Field {
	fieldTagValue := e.GetFieldTagValue(key)
	fieldTagValue.AddValue(value)
	return e
}

// RemoveFieldTagValue remove a tag value
func (e *Field) RemoveFieldTagValue(key, value string) *Field {
	fieldTagValue := e.GetFieldTagValue(key)
	fieldTagValue.RemoveValue(value)
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
		// sort keys
		keys := maps.Keys(e.FieldTags)
		sort.Strings(keys)

		buf.WriteString("`")
		cnt := 0
		for _, k := range keys {
			fieldTagValue := e.FieldTags[k]
			if fieldTagValue != nil && !fieldTagValue.IsEmpty() {
				if cnt != 0 {
					buf.WriteString(" ")
				}
				cnt++
				buf.WriteString(fmt.Sprintf(`%v:"%v"`, k, strings.Join(fieldTagValue.Value, fieldTagValue.Sep)))
			}
		}
		buf.WriteString("`")
	}

	// field comment
	comment := strings.ReplaceAll(e.FieldComment, "\n", ",")
	if comment != "" {
		buf.WriteString("// " + e.FieldComment)
	}
	return buf.String()
}
