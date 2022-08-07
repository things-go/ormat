package deploy

// 运行工作模式,布署
const (
	Local = "local" // 本地
	Dev   = "dev"   // 开发
	Debug = "debug" // 调试
	Uat   = "uat"   // 预发布
	Prod  = "prod"  // 生产
)

var deploy = Prod

// Set 设置布署模式
func Set(m string) {
	switch m {
	case Local, Dev, Debug, Uat, Prod:
		deploy = m
	default:
		deploy = Prod
	}
}

// Get 获取当前的布署模式
func Get() string { return deploy }

// IsTest 测试: 本地,开发或者调试
func IsTest() bool { return IsLocal() || IsDev() || IsDebug() }

// IsRelease 预发或者生产环境
func IsRelease() bool { return IsProd() || IsUat() }

// IsLocal 是否本地模式
func IsLocal() bool { return deploy == Local }

// IsDev 是否开发模式
func IsDev() bool { return deploy == Dev }

// IsDebug 是否调试模式
func IsDebug() bool { return deploy == Debug }

// IsUat 是否预发布模式
func IsUat() bool { return deploy == Uat }

// IsProd 是否生产模式
func IsProd() bool { return deploy == Prod }
