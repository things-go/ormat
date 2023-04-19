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
	FieldTag     string // field tag merge from FieldTags
	FieldComment string // field comment
	IsNullable   bool   // field is null or not
	IsTimestamp  bool   // field Go Type is time.Time
	ColumnGoType string // field column go standard type, e.g. int, uint, float64 ...
	ColumnName   string // field column name in database.
	Type         string // field defined SQL, e.g. varchar(255) NOT NULL DEFAULT ''
	IsSkipColumn bool   // skip filed used for output column
	AssistType   string // field AssistType, e.g. Bool, Uint32, Uint, Int
}
