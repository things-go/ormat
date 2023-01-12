package ast

// Field of a struct
//
//	  FieldName  FieldType  FieldTags                               FieldComment
//
//		 |        |            |                                         |
//	     v        v            v                                         v
//		Foo      int      `json:"foo,omitempty yaml:"foo,omitempty"` // 我是一个注释
//
// FieldTag
//
//	key   v1 sep v2 ...
//	 |     | |
//	 v     v v
//	json:"foo,omitempty"
type Field struct {
	FieldName    string // field name
	FieldType    string // field type
	FieldComment string // field comment
	FieldTag     string // field tag merge from FieldTags
	IsNullable   bool   // field is null or not
	IsTimestamp  bool   // field Go Type is time.Time
	ColumnGoType string // field column go type
	ColumnName   string // field column name in database
	IsSkipColumn bool   // skip filed use for output column
}
