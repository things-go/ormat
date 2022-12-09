package cmd

import (
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/things-go/log"

	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/pkg/consts"
	"github.com/things-go/ormat/pkg/database"
	"github.com/things-go/ormat/pkg/tpl"
	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/runtime"
	"github.com/things-go/ormat/view"
	"github.com/things-go/ormat/view/ast"
)

var genCmd = &cobra.Command{
	Use:     "gen",
	Short:   "Generate model/proto from database",
	Example: "ormat gen",
	RunE: func(*cobra.Command, []string) error {
		c := config.Global
		err := c.Load()
		if err != nil {
			return err
		}
		setupBase(c)
		rt, err := runtime.NewRuntime(c)
		if err != nil {
			return err
		}
		defer database.Close(rt.DB)

		vw := view.New(GetViewModel(rt), c.View)

		mergeProtoEnumFile := ast.ProtobufEnumFile{
			Version: consts.Version,
			Package: vw.Protobuf.Package,
			Options: vw.Protobuf.Options,
			Enums:   make([]*ast.ProtobufEnum, 0, 64),
		}

		list, err := vw.GetDbFile(utils.GetPkgName(c.OutDir))
		if err != nil {
			return err
		}
		for _, v := range list {
			modelFilename := c.OutDir + "/" + v.Filename + ".go"
			_ = utils.WriteFileWithTemplate(modelFilename, tpl.ModelTpl, v)

			cmd, _ := exec.Command("goimports", "-l", "-w", modelFilename).Output()
			_, _ = exec.Command("gofmt", "-l", "-w", modelFilename).Output()
			log.Info("👉 " + strings.TrimSuffix(string(cmd), "\n"))

			if vw.Protobuf.Enabled {
				if enums := v.GetProtobufEnums(); len(enums) > 0 {
					if vw.Protobuf.Merge {
						mergeProtoEnumFile.Enums = append(mergeProtoEnumFile.Enums, enums...)
					} else {
						protoFilename := intoFilename(vw.Protobuf.Dir, v.Filename, ".proto")
						_ = utils.WriteFileWithTemplate(protoFilename, tpl.ProtobufEnumTpl, ast.ProtobufEnumFile{
							Version: consts.Version,
							Package: vw.Protobuf.Package,
							Options: vw.Protobuf.Options,
							Enums:   enums,
						})
						log.Info("👆 " + protoFilename)
					}
				}
			}
		}

		if vw.Protobuf.Enabled && vw.Protobuf.Merge && len(mergeProtoEnumFile.Enums) > 0 {
			enumFilename := intoFilename(vw.Protobuf.Dir, vw.Protobuf.GetMergeFilename(), ".proto")
			_ = utils.WriteFileWithTemplate(enumFilename, tpl.ProtobufEnumTpl, mergeProtoEnumFile)
			log.Info("👆 " + enumFilename)
		}

		log.Info("😄 generate success !!!")
		return nil
	},
}
