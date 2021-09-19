package env

// 运行工作模式,布署
const (
	DeployLocal = "local" // 本地
	DeployDev   = "dev"   // 开发
	DeployDebug = "debug" // 调试
	DeployUat   = "uat"   // 预发布
	DeployProd  = "prod"  // 生产
)

var deploy = DeployProd

// SetDeploy 设置布署模式
func SetDeploy(m string) {
	switch m {
	case DeployLocal, DeployDev, DeployDebug, DeployUat, DeployProd:
		deploy = m
	default:
		deploy = DeployProd
	}
}

// GetDeploy 获取当前的布署模式
func GetDeploy() string {
	return deploy
}

// IsDeployTest 测试: 本地,开发或者调试
func IsDeployTest() bool {
	return IsDeployLocal() || IsDeployDev() || IsDeployDebug()
}

// IsDeployRelease 预发或者生产环境
func IsDeployRelease() bool {
	return IsDeployProd() || IsDeployUat()
}

// IsDeployLocal 是否本地模式
func IsDeployLocal() bool {
	return deploy == DeployLocal
}

// IsDeployDev 是否开发模式
func IsDeployDev() bool {
	return deploy == DeployDev
}

// IsDeployDebug 是否调试模式
func IsDeployDebug() bool {
	return deploy == DeployDebug
}

// IsDeployUat 是否预发布模式
func IsDeployUat() bool {
	return deploy == DeployUat
}

// IsDeployProd 是否生产模式
func IsDeployProd() bool {
	return deploy == DeployProd
}
