package cmd

import (
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/things-go/log"

	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/pkg/consts"
	"github.com/things-go/ormat/pkg/database"
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
			Version:  consts.Version,
			Package:  vw.Protobuf.Package,
			Options:  vw.Protobuf.Options,
			Enums:    nil,
			Template: ast.ProtobufEnumTpl,
		}

		list, err := vw.GetDbFile(utils.GetPkgName(c.OutDir))
		if err != nil {
			return err
		}
		for _, v := range list {
			modelFilename := c.OutDir + "/" + v.Filename + ".go"
			_ = utils.WriteFile(modelFilename, v.Build())

			cmd, _ := exec.Command("goimports", "-l", "-w", modelFilename).Output()
			_, _ = exec.Command("gofmt", "-l", "-w", modelFilename).Output()
			log.Info("👉 " + strings.TrimSuffix(string(cmd), "\n"))

			if c.View.IsOutSQL {
				_ = utils.WriteFile(c.OutDir+"/"+v.Filename+".sql", v.BuildSQL())
			}

			if vw.Protobuf.Enabled {
				if enums := v.GetEnums(); len(enums) > 0 {
					if vw.Protobuf.Merge {
						mergeProtoEnumFile.Enums = append(mergeProtoEnumFile.Enums, enums...)
					} else {
						protoEnumFile := ast.ProtobufEnumFile{
							Version:  consts.Version,
							Package:  vw.Protobuf.Package,
							Options:  vw.Protobuf.Options,
							Enums:    enums,
							Template: ast.ProtobufEnumTpl,
						}
						content := protoEnumFile.Build()
						protoFilename := intoFilename(vw.Protobuf.Dir, v.Filename, ".proto")
						_ = utils.WriteFile(protoFilename, content)
						log.Info("👆 " + protoFilename)

					}
				}
			}
		}

		if vw.Protobuf.Enabled &&
			vw.Protobuf.Merge &&
			len(mergeProtoEnumFile.Enums) > 0 {
			mergeFilename := vw.Protobuf.GetMergeFilename()
			enumFilename := intoFilename(vw.Protobuf.Dir, mergeFilename, ".proto")
			content := mergeProtoEnumFile.Build()
			_ = utils.WriteFile(enumFilename, content)
			log.Info("👆 " + enumFilename)
		}

		log.Info("😄 generate success !!!")
		return nil
	},
}
