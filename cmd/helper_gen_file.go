package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/alecthomas/chroma/quick"
	"github.com/things-go/log"

	"github.com/things-go/ormat/pkg/consts"
	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/view/ast"
)

type generateFile struct {
	Files     []*ast.File
	Template  *template.Template
	OutputDir string
	// only use for enum
	Merge         bool
	MergeFilename string
	Package       string
	Options       map[string]string
	Suffix        string
	GenFunc       func(filename string, t *template.Template, data any)
}

func (g *generateFile) runGenModel() {
	for _, v := range g.Files {
		g.GenFunc(
			intoFilename(g.OutputDir, v.Filename, ".go"),
			g.Template,
			v,
		)
	}
	log.Info("😄 generate success !!!")
}

func (g *generateFile) runGenEnum() {
	packageName := utils.GetPkgName(g.OutputDir)
	mergeProtoEnumFile := ast.ProtobufEnumFile{
		Version:     consts.Version,
		PackageName: packageName,
		Package:     g.Package,
		Options:     g.Options,
		Enums:       make([]*ast.ProtobufEnum, 0, 64),
	}
	for _, v := range g.Files {
		if enums := v.GetProtobufEnums(); len(enums) > 0 {
			if g.Merge {
				mergeProtoEnumFile.Enums = append(mergeProtoEnumFile.Enums, enums...)
			} else {
				g.GenFunc(
					intoFilename(g.OutputDir, v.Filename, g.Suffix),
					g.Template,
					ast.ProtobufEnumFile{
						Version:     consts.Version,
						PackageName: packageName,
						Package:     g.Package,
						Options:     g.Options,
						Enums:       enums,
					},
				)
			}
		}
	}
	if g.Merge && len(mergeProtoEnumFile.Enums) > 0 {
		if g.MergeFilename == "" {
			g.MergeFilename = utils.GetPkgName(g.OutputDir)
		}
		g.GenFunc(
			intoFilename(g.OutputDir, g.MergeFilename, g.Suffix),
			g.Template,
			mergeProtoEnumFile,
		)
	}

	log.Info("😄 generate success !!!")
}

func genModelFile(filename string, t *template.Template, data any) {
	_ = utils.WriteFileWithTemplate(filename, t, data)

	cmd, _ := exec.Command("goimports", "-l", "-w", filename).Output()
	_, _ = exec.Command("gofmt", "-l", "-w", filename).Output()
	log.Info("👉 " + strings.TrimSuffix(string(cmd), "\n"))
}

func genEnumFile(filename string, t *template.Template, data any) {
	_ = utils.WriteFileWithTemplate(filename, t, data)
	log.Info("👆 " + filename)
}

func showInfo(filename string, t *template.Template, data any) {
	b, _ := json.MarshalIndent(data, " ", "  ")
	quick.Highlight(os.Stdout, string(b), "JSON", "terminal", "solarized-dark")
}
