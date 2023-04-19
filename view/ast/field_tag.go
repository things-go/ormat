package ast

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

// FieldTagValues field tag values
type FieldTagValues struct {
	Separate string   // tag separate like [",", ";", " "]
	Value    []string // tag value
}

// NewFiledTagValues new FieldTagValues instance with default separate `,`.
func NewFiledTagValues() *FieldTagValues {
	return &FieldTagValues{Separate: ","}
}

// IsZero value number is zero
func (f *FieldTagValues) IsZero() bool { return len(f.Value) == 0 }

// SetSeparate set separate.
func (f *FieldTagValues) SetSeparate(separate string) *FieldTagValues {
	f.Separate = separate
	return f
}

// AddValue add value
func (f *FieldTagValues) AddValue(value string) *FieldTagValues {
	f.Value = append(f.Value, value)
	return f
}

// RemoveValue remove value
func (f *FieldTagValues) RemoveValue(value string) *FieldTagValues {
	tagsValues := f.Value
	for i, vv := range tagsValues {
		if vv == value {
			f.Value = append(tagsValues[:i], tagsValues[i+1:]...)
		}
	}
	return f
}

// FieldTags contains multiple tags
//
//	key   v1 sep v2 ...
//	 |     | |
//	 v     v v
//	json:"foo,omitempty"
type FieldTags struct {
	inner map[string]*FieldTagValues // field tags
}

func NewFieldTags() *FieldTags {
	return &FieldTags{
		inner: make(map[string]*FieldTagValues),
	}
}

// Get get a tag.
func (f *FieldTags) Get(key string) *FieldTagValues {
	fieldTagValue, ok := f.inner[key]
	if !ok {
		fieldTagValue = NewFiledTagValues()
		f.inner[key] = fieldTagValue
	}
	return fieldTagValue
}

// Add add a tag.
func (f *FieldTags) Add(key string, tagValue *FieldTagValues) *FieldTags {
	if tagValue != nil {
		f.inner[key] = tagValue
	}
	return f
}

// Remove remove a tag.
func (f *FieldTags) Remove(key string) *FieldTags {
	delete(f.inner, key)
	return f
}

// SetTagSeparate set the tag separate
func (f *FieldTags) SetTagSeparate(key, sep string) *FieldTags {
	fieldTagValue := f.Get(key)
	fieldTagValue.SetSeparate(sep)
	return f
}

// AddTagValue add a tag value
func (f *FieldTags) AddTagValue(key, value string) *FieldTags {
	fieldTagValue := f.Get(key)
	fieldTagValue.AddValue(value)
	return f
}

// RemoveTagValue remove a tag value
func (f *FieldTags) RemoveTagValue(key, value string) *FieldTags {
	fieldTagValue := f.Get(key)
	fieldTagValue.RemoveValue(value)
	return f
}

func (f *FieldTags) IntoFieldTag() string {
	var buf strings.Builder

	if len(f.inner) > 0 {
		keys := maps.Keys(f.inner)
		sort.Strings(keys) // sort keys

		cnt := 0
		for _, k := range keys {
			if fieldTagValue := f.inner[k]; fieldTagValue != nil && !fieldTagValue.IsZero() {
				if cnt != 0 {
					buf.WriteString(" ")
				}
				cnt++
				buf.WriteString(fmt.Sprintf(`%v:"%v"`, k, strings.Join(fieldTagValue.Value, fieldTagValue.Separate)))
			}
		}
	}
	return buf.String()
}
