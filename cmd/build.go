package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/things-go/log"

	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/runtime"
	"github.com/things-go/ormat/view"
	"github.com/things-go/ormat/view/driver"
)

var inputFile string

func init() {
	buildCmd.Flags().StringVarP(&inputFile, "input", "i", "", "input file")
	buildCmd.Flags().StringVarP(&outDir, "out", "o", "", "out directory")
	buildCmd.MarkFlagRequired("input")
}

var buildCmd = &cobra.Command{
	Use:     "build",
	Short:   "Generate model from sql",
	Example: "ormat build",
	RunE: func(*cobra.Command, []string) error {
		rt, err := runtime.NewRuntime(false)
		if err != nil {
			return err
		}
		c := rt.Config

		outDir := outDir
		if outDir == "" {
			outDir = c.OutDir
		}

		content, err := os.ReadFile(inputFile)
		if err != nil {
			return err
		}
		vw := view.New(&driver.SQL{
			CreateTableSQL:   string(content),
			CustomDefineType: c.TypeDefine,
		}, c.View)
		list, err := vw.GetDbFile(utils.GetPkgName(outDir))
		if err != nil {
			return err
		}
		for _, v := range list {
			path := outDir + "/" + v.Filename + ".go"
			_ = utils.WriteFile(path, []byte(v.Build()))

			cmd, _ := exec.Command("goimports", "-l", "-w", path).Output()
			log.Info(strings.TrimSuffix(string(cmd), "\n"))
			_, _ = exec.Command("gofmt", "-l", "-w", path).Output()

			if c.View.IsOutSQL {
				_ = utils.WriteFile(outDir+"/"+v.Filename+".sql", []byte(v.BuildSQL()))
			}
		}
		log.Info("generate success !!!")
		return nil
	},
}
