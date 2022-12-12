package ast

import (
	"strings"
)

// ImportsHeads import head options
var ImportsHeads = map[string]string{
	"time.Time":      `"time"`,
	"gorm.Model":     `"gorm.io/gorm"`,
	"fmt":            `"fmt"`,
	"datatypes.JSON": `"gorm.io/datatypes"`,
	"datatypes.Date": `"gorm.io/datatypes"`,
}

// IntoAbbrTableName 获取表名缩写
func IntoAbbrTableName(tableName string) string {
	ss := strings.Split(tableName, "_")
	tableName = ""
	for _, vv := range ss {
		if len(vv) > 0 {
			tableName += string(vv[0])
		}
	}
	return tableName
}

// IntoComment 转换注释
func IntoComment(comment, defaultComment, old, new string) string {
	if comment == "" {
		comment = defaultComment
	} else {
		comment = strings.ReplaceAll(strings.TrimSpace(comment), old, new)
	}
	return comment
}
