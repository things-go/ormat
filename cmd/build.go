package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/things-go/log"

	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/view"
	"github.com/things-go/ormat/view/driver"
)

var inputFile string
var protobufOptions []string

func init() {
	buildCmd.Flags().StringVarP(&inputFile, "input", "i", "", "input file")
	buildCmd.Flags().StringP("out", "o", "", "model out directory")
	buildCmd.Flags().StringP("dir", "d", "", "protobuf out directory")
	buildCmd.Flags().BoolP("enabled", "e", true, "enabled generate protobuf")
	buildCmd.Flags().StringP("package", "p", "", "protobuf package name")
	buildCmd.Flags().StringSliceVarP(&protobufOptions, "options", "t", nil, "protobuf options key value")

	// buildCmd.MarkFlagsRequiredTogether(
	// 	"enabled",
	// 	"dir",
	// 	"package",
	// 	"options",
	// )
	buildCmd.MarkFlagRequired("input")
	// buildCmd.MarkFlagRequired("out")

	viper.BindPFlag("outDir", buildCmd.Flags().Lookup("out"))
	// viper.BindPFlag("view.protobuf.enabled", buildCmd.Flags().Lookup("enabled"))
	// viper.BindPFlag("view.protobuf.dir", buildCmd.Flags().Lookup("dir"))
	// viper.BindPFlag("view.protobuf.package", buildCmd.Flags().Lookup("package"))
}

var buildCmd = &cobra.Command{
	Use:     "build",
	Short:   "Generate model from sql",
	Example: "ormat build",
	RunE: func(*cobra.Command, []string) error {
		c := config.Global
		err := c.Load()
		if err != nil {
			return err
		}

		outDir := c.OutDir
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
			_ = utils.WriteFile(path, v.Build())

			cmd, _ := exec.Command("goimports", "-l", "-w", path).Output()
			log.Info("👉 " + strings.TrimSuffix(string(cmd), "\n"))
			_, _ = exec.Command("gofmt", "-l", "-w", path).Output()

			if c.View.IsOutSQL {
				_ = utils.WriteFile(outDir+"/"+v.Filename+".sql", v.BuildSQL())
			}
		}
		log.Info("😄 generate success !!!")
		return nil
	},
}
