package cmd

import (
	"os/exec"
	"strings"
	"text/template"

	"github.com/things-go/log"

	"github.com/things-go/ormat/pkg/consts"
	"github.com/things-go/ormat/pkg/tpl"
	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/view"
	"github.com/things-go/ormat/view/ast"
)

func genModelFile(astFiles []*ast.File, protobufConfig *view.Protobuf, outDir string) {
	mergeProtoEnumFile := &ast.ProtobufEnumFile{
		Version: consts.Version,
		Package: protobufConfig.Package,
		Options: protobufConfig.Options,
		Enums:   make([]*ast.ProtobufEnum, 0, 64),
	}
	for _, v := range astFiles {
		modelFilename := outDir + "/" + v.Filename + ".go"
		_ = utils.WriteFileWithTemplate(modelFilename, tpl.ModelTpl, v)

		cmd, _ := exec.Command("goimports", "-l", "-w", modelFilename).Output()
		_, _ = exec.Command("gofmt", "-l", "-w", modelFilename).Output()
		log.Info("👉 " + strings.TrimSuffix(string(cmd), "\n"))

		if protobufConfig.Enabled {
			if enums := v.GetProtobufEnums(); len(enums) > 0 {
				if protobufConfig.Merge {
					mergeProtoEnumFile.Enums = append(mergeProtoEnumFile.Enums, enums...)
				} else {
					enumFilename := intoFilename(protobufConfig.Dir, v.Filename, ".proto")
					_ = utils.WriteFileWithTemplate(enumFilename, tpl.ProtobufEnumTpl, ast.ProtobufEnumFile{
						Version: consts.Version,
						Package: protobufConfig.Package,
						Options: protobufConfig.Options,
						Enums:   enums,
					})
					log.Info("👆 " + enumFilename)
				}
			}
		}
	}
	if protobufConfig.Enabled && protobufConfig.Merge && len(mergeProtoEnumFile.Enums) > 0 {
		enumFilename := intoFilename(protobufConfig.Dir, protobufConfig.GetMergeFilename(), ".proto")
		_ = utils.WriteFileWithTemplate(enumFilename, tpl.ProtobufEnumTpl, mergeProtoEnumFile)
		log.Info("👆 " + enumFilename)
	}

	log.Info("😄 generate success !!!")
}

func genEnumFile(astFiles []*ast.File, protobufConfig *view.Protobuf, usedTemplate *template.Template) {
	mergeProtoEnumFile := ast.ProtobufEnumFile{
		Version: consts.Version,
		Package: protobufConfig.Package,
		Options: protobufConfig.Options,
		Enums:   nil,
	}

	for _, v := range astFiles {
		if enums := v.GetProtobufEnums(); len(enums) > 0 {
			if protobufConfig.Merge {
				mergeProtoEnumFile.Enums = append(mergeProtoEnumFile.Enums, enums...)
			} else {
				enumFilename := intoFilename(protobufConfig.Dir, v.Filename, protobufConfig.Suffix)
				_ = utils.WriteFileWithTemplate(enumFilename, usedTemplate, ast.ProtobufEnumFile{
					Version: consts.Version,
					Package: protobufConfig.Package,
					Options: protobufConfig.Options,
					Enums:   enums,
				})
				log.Info("👆 " + enumFilename)
			}
		}
	}
	if protobufConfig.Merge && len(mergeProtoEnumFile.Enums) > 0 {
		mergeFilename := protobufConfig.GetMergeFilename()
		enumFilename := intoFilename(protobufConfig.Dir, mergeFilename, protobufConfig.Suffix)
		_ = utils.WriteFileWithTemplate(enumFilename, usedTemplate, mergeProtoEnumFile)
		log.Info("👆 " + enumFilename)
	}

	log.Info("😄 generate success !!!")
}
