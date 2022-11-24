package cmd

import (
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/things-go/log"

	"github.com/things-go/ormat/pkg/database"
	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/runtime"
	"github.com/things-go/ormat/view"
)

var genCmd = &cobra.Command{
	Use:     "gen",
	Short:   "Generate model from database",
	Example: "ormat gen",
	RunE: func(*cobra.Command, []string) error {
		rt, err := runtime.NewRuntime(true)
		if err != nil {
			return err
		}
		defer database.Close(rt.DB)

		c := rt.Config
		vw := view.New(GetViewModel(rt), c.View)

		list, err := vw.GetDbFile(utils.GetPkgName(c.OutDir))
		if err != nil {
			return err
		}
		for _, v := range list {
			modelFilename := c.OutDir + "/" + v.Filename + ".go"
			_ = utils.WriteFile(modelFilename, v.Build())

			cmd, _ := exec.Command("goimports", "-l", "-w", modelFilename).Output()
			log.Info("👉 " + strings.TrimSuffix(string(cmd), "\n"))
			_, _ = exec.Command("gofmt", "-l", "-w", modelFilename).Output()

			if c.View.IsOutSQL {
				_ = utils.WriteFile(c.OutDir+"/"+v.Filename+".sql", v.BuildSQL())
			}

			if vw.Protobuf.Enabled {
				content := v.BuildProtobufEnum()
				if len(content) > 0 {
					protoFilename := vw.Protobuf.Dir + "/" + v.Filename + ".proto"
					_ = utils.WriteFile(protoFilename, content)
					log.Info("👆 " + protoFilename)
				}
			}
		}
		log.Info("😄 generate success !!!")
		return nil
	},
}
