package builder

import (
	"os"
	"runtime"
	"text/template"
)

var (
	// Model 型号 由外部ldflags指定
	Model = "unknown"
	// Name 应用名称 由外部ldflags指定
	Name = "unknown"
	// Version 版本 由外部ldflags指定
	Version = "unknown"
	// GitCommit git提交版本(短) 由外部ldflags指定
	GitCommit = "unknown"
	// GitFullCommit git提交版本(完整) 由外部ldflags指定
	GitFullCommit = "unknown"
	// GitFullCommit git标签 由外部ldflags指定
	GitTag = "unknown"
	// BuildTime 编译日期 由外部ldflags指定
	BuildTime = "unknown"
)

const versionTpl = `  Name:             {{.Name}}
  Model:            {{.Model}}
  Version:          {{.Version}}
  Git commit:       {{.GitCommit}}
  Git full commit:  {{.GitFullCommit}}
  Git tag:          {{.GitTag}}
  Build time:       {{.BuildTime}}
  Go version:       {{.GoVersion}}
  OS/Arch:          {{.GOOS}}/{{.GOARCH}}
  NumCPU:           {{.NumCPU}}
`

// Version 版本信息
type Ver struct {
	Name          string
	Model         string
	Version       string
	GitCommit     string
	GitFullCommit string
	GitTag        string
	BuildTime     string
	GoVersion     string
	GOOS          string
	GOARCH        string
	NumCPU        int
}

// PrintVersion 打印版本信息至os.Stdout
func PrintVersion() {
	v := Ver{
		Name,
		Model,
		Version,
		GitCommit,
		GitFullCommit,
		GitTag,
		BuildTime,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
	}
	template.Must(template.New("version").Parse(versionTpl)).
		Execute(os.Stdout, v) // nolint: errcheck
}
