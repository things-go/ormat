package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// WorkDir 获取当前工作目录
func WorkDir() string {
	dir, _ := os.Getwd()
	return strings.ReplaceAll(dir, "\\", "/")
}

func ExecutableDir() string {
	dir, _ := os.Executable()
	exPath := filepath.Dir(dir)
	return strings.ReplaceAll(exPath, "\\", "/")
}

func GetPkgName(path string) string {
	pkgName := filepath.Base(path)
	if pkgName == "" || pkgName == "." {
		path, _ = os.Getwd()
		_, pkgName = filepath.Split(path)
	}
	return pkgName
}
