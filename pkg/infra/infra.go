package infra

import (
	"os"
	"path/filepath"
	"strings"
)

const Version = "v0.0.1-rc10"

// GetWd 获取当前工作目录
func GetWd() string {
	dir, _ := os.Getwd()
	return strings.ReplaceAll(dir, "\\", "/")
}

func GetExecutableDir() string {
	dir, _ := os.Executable()
	exPath := filepath.Dir(dir)
	return strings.ReplaceAll(exPath, "\\", "/")
}

func GetPkgName(path string) string {
	_, pkgName := filepath.Split(path)
	if pkgName == "" || pkgName == "." {
		path, _ = os.Getwd()
		_, pkgName = filepath.Split(path)
	}
	return pkgName
}
