package ast

import (
	"regexp"
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

// t.Logf("%#v", rEnum.FindStringSubmatch(` 11 [@enum:{"0":["none","空","空注释"],"1":["key1","键1","键1注释"],"2":["key2","键2","3":["key3","键3"]]}] 11k l23123 人11`))
var rEnum = regexp.MustCompile(`^.*?\[@.*?(?i:(?:enum|status)+):\s*(.*)\].*?`)

// MatchEnumAnnotation 匹配枚举注解
func MatchEnumAnnotation(comment string) string {
	match := rEnum.FindStringSubmatch(comment)
	if len(match) == 2 {
		return strings.TrimSpace(match[1])
	}
	return ""
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
