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

func GetPkgName(path string) string {
	pkgName := filepath.Base(path)
	if pkgName == "" || pkgName == "." {
		workdir := WorkDir()
		pkgName = filepath.Base(workdir)
	}
	return pkgName
}
