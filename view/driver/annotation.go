package driver

import (
	"regexp"
)

var rAutoIncrement = regexp.MustCompile(` (AUTO_INCREMENT=\d+){1} `)
