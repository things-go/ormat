package matcher

import (
	"regexp"
	"strings"
)

var rJSONTag = regexp.MustCompile(`^.*\[@(?i:jsontag):\s*([^\[\]]*)\].*`)
var rAffixJSONTag = regexp.MustCompile(`^.*\[@(affix)\s*\].*`)

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
