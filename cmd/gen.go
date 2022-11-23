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
	Short:   "generate model from sql",
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
			path := c.OutDir + "/" + v.Filename + ".go"
			_ = utils.WriteFile(path, []byte(v.Build()))

			cmd, _ := exec.Command("goimports", "-l", "-w", path).Output()
			log.Info(strings.TrimSuffix(string(cmd), "\n"))
			_, _ = exec.Command("gofmt", "-l", "-w", path).Output()

			if c.View.IsOutSQL {
				_ = utils.WriteFile(c.OutDir+"/"+v.Filename+".sql", []byte(v.BuildSQL()))
			}
		}
		log.Info("generate success !!!")
		return nil
	},
}
