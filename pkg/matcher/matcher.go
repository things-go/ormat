package matcher

import (
	"regexp"
	"strings"
)

var rEnum = regexp.MustCompile(`^.*\[@(?:enum|status):\s*({.*})\s*\].*`)
var rJSONTag = regexp.MustCompile(`^.*\[@(?i:jsontag):\s*([^\[\]]*)\].*`)
var rAffixJSONTag = regexp.MustCompile(`^.*\[@(affix)\s*\].*`)

// EnumAnnotation 匹配枚举注解
// [@enum:{...}]
// [@status:{...}]
func EnumAnnotation(comment string) string {
	match := rEnum.FindStringSubmatch(comment)
	if len(match) == 2 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

// JsonTag 匹配json标签
// [@jsontag:id,omitempty]
func JsonTag(comment string) string {
	match := rJSONTag.FindStringSubmatch(comment)
	if len(match) == 2 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

// HasAffixJSONTag 是否有 affix, 增加 json 标签 `,string`
// [@affix]
func HasAffixJSONTag(comment string) bool {
	match := rAffixJSONTag.FindStringSubmatch(comment)
	return len(match) == 2 && strings.TrimSpace(match[1]) == "affix"
}
