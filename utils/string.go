package utils

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"
)

var (
	// commonInitialisms is a set of common initialisms.
	// source: https://github.com/golang/lint/blob/master/lint.go
	commonInitialisms = map[string]struct{}{
		"ACL":   {},
		"API":   {},
		"ASCII": {},
		"CPU":   {},
		"CSS":   {},
		"DNS":   {},
		"EOF":   {},
		"GUID":  {},
		"HTML":  {},
		"HTTP":  {},
		"HTTPS": {},
		"ID":    {},
		"IP":    {},
		"JSON":  {},
		"LHS":   {},
		"QPS":   {},
		"RAM":   {},
		"RHS":   {},
		"RPC":   {},
		"SLA":   {},
		"SMTP":  {},
		"SQL":   {},
		"SSH":   {},
		"TCP":   {},
		"TLS":   {},
		"TTL":   {},
		"UDP":   {},
		"UI":    {},
		"UID":   {},
		"UUID":  {},
		"URI":   {},
		"URL":   {},
		"UTF8":  {},
		"VM":    {},
		"XML":   {},
		"XMPP":  {},
		"XSRF":  {},
		"XSS":   {},
	}
	defaultReplacer *strings.Replacer
)

func init() {
	initialismForReplacer := make([]string, 0, len(commonInitialisms)*2)
	for s := range commonInitialisms {
		initialismForReplacer = append(initialismForReplacer, s, strings.Title(strings.ToLower(s)))
	}

	defaultReplacer = strings.NewReplacer(initialismForReplacer...)
}

// Recombine 转换驼峰字符串为用delimiter分隔的字符串, 特殊字符由DefaultInitialisms决定取代
// example: delimiter = '_'
// 空字符 -> 空字符
// HelloWorld -> hello_world
// Hello_World -> hello_world
// HiHello_World -> hi_hello_world
// IDCom -> id_com
// IDcom -> idcom
// nameIDCom -> name_id_com
// nameIDcom -> name_idcom
func Recombine(str string, delimiter byte, enableLint bool) string {
	str = strings.TrimSpace(str)
	if str == "" {
		return ""
	}
	if enableLint {
		str = defaultReplacer.Replace(str)
	}

	var isLastCaseUpper bool
	var isCurrCaseUpper bool
	var isNextCaseUpper bool
	var isNextNumberUpper bool
	var buf = strings.Builder{}

	for i, v := range str[:len(str)-1] {
		isNextCaseUpper = str[i+1] >= 'A' && str[i+1] <= 'Z'
		isNextNumberUpper = str[i+1] >= '0' && str[i+1] <= '9'

		if i > 0 {
			if isCurrCaseUpper {
				if isLastCaseUpper && (isNextCaseUpper || isNextNumberUpper) {
					buf.WriteRune(v)
				} else {
					if str[i-1] != delimiter && str[i+1] != delimiter {
						buf.WriteRune(rune(delimiter))
					}
					buf.WriteRune(v)
				}
			} else {
				buf.WriteRune(v)
				if i == len(str)-2 && (isNextCaseUpper && !isNextNumberUpper) {
					buf.WriteRune(rune(delimiter))
				}
			}
		} else {
			isCurrCaseUpper = true
			buf.WriteRune(v)
		}
		isLastCaseUpper = isCurrCaseUpper
		isCurrCaseUpper = isNextCaseUpper
	}

	buf.WriteByte(str[len(str)-1])

	return strings.ToLower(buf.String())
}

// UnRecombine 转换sep分隔的字符串为驼峰字符串
// example: delimiter = '_'
// 空字符 -> 空字符
// hello_world -> HelloWorld
func UnRecombine(str string, delimiter byte, enableLint bool) string {
	str = strings.TrimSpace(str)
	if str == "" {
		return ""
	}

	var b strings.Builder
	var words []string

	for i, s := 0, str; s != ""; s = s[i:] { // split on upper letter or _
		i = strings.IndexFunc(s[1:], unicode.IsUpper) + 1
		if i <= 0 {
			i = len(s)
		}
		word := s[:i]
		words = append(words, strings.Split(word, string(delimiter))...)
	}

	for i, word := range words {
		if enableLint {
			u := strings.ToUpper(word)
			if _, ok := commonInitialisms[u]; ok {
				b.WriteString(u)
				continue
			}
		}

		word = removeInvalidChars(word, i == 0) // on 0 remove first digits
		if word == "" {
			continue
		}

		out := strings.ToUpper(string(word[0]))
		if len(word) > 1 {
			out += strings.ToLower(word[1:])
		}
		b.WriteString(out)
	}

	if b.Len() == 0 { // check if this is number
		if _, err := strconv.Atoi(str); err == nil {
			b.WriteString("Key")
			b.WriteString(str)
		}
	}

	return b.String()
}

// SnakeCase 转换驼峰字符串为用'_'分隔的字符串,特殊字符由DefaultInitialisms决定取代
// example2: delimiter = '_' initialisms = DefaultInitialisms
// IDCom -> id_com
// IDcom -> idcom
// nameIDCom -> name_id_com
// nameIDcom -> name_idcom
func SnakeCase(str string, enableLint bool) string {
	return Recombine(str, '_', enableLint)
}

// Kebab 转换驼峰字符串为用'-'分隔的字符串,特殊字符由DefaultInitialisms决定取代
// example2: delimiter = '-' initialisms = DefaultInitialisms
// IDCom -> id-com
// IDcom -> idcom
// nameIDCom -> name-id-com
// nameIDcom -> name-idcom
func Kebab(str string, enableLint bool) string {
	return Recombine(str, '-', enableLint)
}

// CamelCase to camel case string
// id_com -> IDCom
// idcom -> Idcom
// name_id_com -> NameIDCom
// name_idcom -> NameIdcom
func CamelCase(str string, enableLint bool) string {
	return UnRecombine(str, '_', enableLint)
}

// SmallCamelCase to small camel case string
// id_com -> idCom
// idcom -> idcom
// name_id_com -> nameIDCom
// name_idcom -> nameIdcom
func SmallCamelCase(fieldName string, enableLint bool) string {
	fieldName = CamelCase(fieldName, enableLint)
	if enableLint {
		for k := range commonInitialisms {
			if strings.HasPrefix(fieldName, k) {
				return strings.Replace(fieldName, k, strings.ToLower(k), 1)
			}
		}
	}
	return LowTitle(fieldName)
}

func removeInvalidChars(s string, removeFirstDigit bool) string {
	var buf bytes.Buffer

	for _, b := range []byte(s) {
		if b >= 97 && b <= 122 { // a-z
			buf.WriteByte(b)
			continue
		}
		if b >= 65 && b <= 90 { // A-Z
			buf.WriteByte(b)
			continue
		}
		if b >= 48 && b <= 57 { // 0-9
			if !removeFirstDigit || buf.Len() > 0 {
				buf.WriteByte(b)
				continue
			}
		}
	}

	return buf.String()
}

// isSeparator reports whether the rune could mark a word boundary.
// TODO: update when package unicode captures more of the properties.
// see strings isSeparator
func isSeparator(r rune) bool {
	// ASCII alphanumerics and underscore are not separators
	if r <= 0x7F {
		switch {
		case r >= '0' && r <= '9':
			return false
		case r >= 'a' && r <= 'z':
			return false
		case r >= 'A' && r <= 'Z':
			return false
		case r == '_':
			return false
		}
		return true
	}

	// Letters and digits are not separators
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	// Otherwise, all we can do for now is treat spaces as separators.
	return unicode.IsSpace(r)
}

// LowTitle 首字母小写
// see strings.Title
func LowTitle(s string) string {
	// Use a closure here to remember state.
	// Hackish but effective. Depends on Map scanning in order and calling
	// the closure once per rune.
	prev := ' '
	return strings.Map(func(r rune) rune {
		if isSeparator(prev) {
			prev = r
			return unicode.ToLower(r)
		}
		prev = r
		return r
	}, s)
}
