package utils

import (
	"strings"
	"unicode"
)

// SplitCase 转换驼峰字符串为用delimiter分隔的字符串
// example: delimiter = '_'
// 空字符 -> 空字符
// HelloWorld -> hello_world
// Hello_World -> hello_world
// HiHello_World -> hi_hello_world
// IdCom -> id_com
// Idcom -> idcom
// nameIdCom -> name_id_com
// nameIdcom -> name_idcom
func SplitCase(str string, delimiter byte) string {
	str = strings.TrimSpace(str)
	if str == "" {
		return ""
	}
	var isPrevSpecial bool

	t := make([]byte, 0, len(str)+8)
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for i := 0; i < len(str); i++ {
		c := str[i]
		isCureNumber := isASCIIDigit(c)
		if isCureNumber {
			if i == 0 {
				t = append(t, 'x', delimiter)
			}
		} else {
			if isASCIIUpper(c) {
				c += 'a' - 'A'
				if i > 0 && !isPrevSpecial {
					t = append(t, delimiter)
				}
			}
		}
		isPrevSpecial = isCureNumber || c == delimiter
		t = append(t, c) // Guaranteed not lower case.
	}
	return string(t)
}

func JoinCase(s string, delimiter byte) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == delimiter {
		// Need a capital letter; drop the delimiter.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == delimiter && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// SnakeCase 转换驼峰字符串为用'_'分隔的字符串
// example2: delimiter = '_' initialisms = DefaultInitialisms
// IdCom -> id_com
// Idcom -> idcom
// nameIdCom -> name_id_com
// nameIdcom -> name_idcom
func SnakeCase(str string) string {
	return SplitCase(str, '_')
}

// Kebab 转换驼峰字符串为用'-'分隔的字符串
// example2: delimiter = '-' initialisms = DefaultInitialisms
// IdCom -> id-com
// Idcom -> idcom
// nameIdCom -> name-id-com
// nameIdcom -> name-idcom
func Kebab(str string) string {
	return SplitCase(str, '-')
}

// SmallCamelCase to small camel case string
// id_com -> idCom
// idcom -> idcom
// name_id_com -> nameIdCom
// name_idcom -> nameIdcom
func SmallCamelCase(fieldName string) string {
	return LowTitle(CamelCase(fieldName))
}

// CamelCase returns the CamelCased name.
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func CamelCase(s string) string {
	return JoinCase(s, '_')
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

// Is c an ASCII upper-case letter?
func isASCIIUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
