//go:generate stringer -type=Deploy -linecomment
package deploy

type Deploy int

const (
	None Deploy = iota // none
	Dev                // dev
	Test               // test
	Uat                // uat
	Prod               // prod
)

var deploy = None

// Convert m to Deploy
func Convert(m string) Deploy {
	switch m {
	case Test.String():
		return Test
	case Dev.String():
		return Dev
	case Uat.String():
		return Uat
	case Prod.String():
		return Prod
	default:
		return None
	}
}

// IsUat 是否预发布模式
func IsUat() bool { return deploy == Uat }

// IsProduction 是否生产模式
func IsProduction() bool { return deploy == Prod }

// IsRelease 发布, 预发或者生产环境
func IsRelease() bool { return IsUat() || IsProduction() }
