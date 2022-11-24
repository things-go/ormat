package ast

import (
	"encoding/json"
	"regexp"
	"strings"
)

// interval
const delimTab = "\t"
const delimLF = "\n"
const delimSpace = " "

// ImportsHeads import head options
var ImportsHeads = map[string]string{
	"string":         `"string"`,
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

// ParseEnumAnnotation 解析枚举注解.
func ParseEnumAnnotation(annotation string) (map[string][]string, error) {
	var mp map[string][]string

	err := json.Unmarshal([]byte(annotation), &mp)
	if err != nil {
		return nil, err
	}
	return mp, nil
}
