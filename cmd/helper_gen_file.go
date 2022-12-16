package cmd

import (
	"os/exec"
	"strings"
	"text/template"

	"github.com/things-go/log"

	"github.com/things-go/ormat/pkg/consts"
	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/view/ast"
)

type generateFile struct {
	Files         []*ast.File
	Template      *template.Template
	OutputDir     string
	Merge         bool
	MergeFilename string
	Package       string
	Options       map[string]string
	Suffix        string
	GenFunc       func(filename string, t *template.Template, data any)
}

func (g *generateFile) runGenModel() {
	packageName := utils.GetPkgName(g.OutputDir)
	mergeFile := ast.File{
		Version:     consts.Version,
		Filename:    g.MergeFilename,
		PackageName: packageName,
		Imports:     make(map[string]struct{}),
		Structs:     make([]*ast.Struct, 0, 512),
	}
	for _, v := range g.Files {
		if g.Merge {
			for k := range v.Imports {
				mergeFile.Imports[k] = struct{}{}
			}
			mergeFile.Structs = append(mergeFile.Structs, v.Structs...)
		} else {
			g.GenFunc(
				intoFilename(g.OutputDir, v.Filename, g.Suffix),
				g.Template,
				v,
			)
		}
	}
	if g.Merge && len(mergeFile.Structs) > 0 {
		if mergeFile.Filename == "" {
			mergeFile.Filename = packageName
		}
		g.GenFunc(
			intoFilename(g.OutputDir, mergeFile.Filename, g.Suffix),
			g.Template,
			mergeFile,
		)
	}

	log.Info("😄 generate success !!!")
}

func (g *generateFile) runGenEnum() {
	packageName := utils.GetPkgName(g.OutputDir)
	mergeEnumFile := ast.ProtobufEnumFile{
		Version:     consts.Version,
		PackageName: packageName,
		Package:     g.Package,
		Options:     g.Options,
		Enums:       make([]*ast.ProtobufEnum, 0, 64),
	}
	for _, v := range g.Files {
		if enums := v.GetProtobufEnums(); len(enums) > 0 {
			if g.Merge {
				mergeEnumFile.Enums = append(mergeEnumFile.Enums, enums...)
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
	if g.Merge && len(mergeEnumFile.Enums) > 0 {
		if g.MergeFilename == "" {
			g.MergeFilename = packageName
		}
		g.GenFunc(
			intoFilename(g.OutputDir, g.MergeFilename, g.Suffix),
			g.Template,
			mergeEnumFile,
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

func showInformation(filename string, t *template.Template, data any) {
	JSON(data)
}
