package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/things-go/log"

	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/view"
	"github.com/things-go/ormat/view/ast"
	"github.com/things-go/ormat/view/driver"
)

var inputFile []string
var outDir string
var protobuf = view.Protobuf{
	Enabled: false,
	Merge:   false,
	Dir:     "",
	Package: "",
	Options: nil,
}

func init() {
	buildCmd.PersistentFlags().StringSliceVarP(&inputFile, "input", "i", nil, "input file")
	buildCmd.PersistentFlags().StringVarP(&outDir, "out", "o", "", "model out directory")

	buildCmd.Flags().BoolVarP(&protobuf.Enabled, "enabled", "e", false, "protobuf enabled or not(default: false)")
	buildCmd.Flags().BoolVarP(&protobuf.Merge, "merge", "m", false, "protobuf merge in a file or not(default: false)")
	buildCmd.Flags().StringVarP(&protobuf.Dir, "dir", "d", "", "protobuf out directory")
	buildCmd.Flags().StringVarP(&protobuf.Package, "package", "p", "", "protobuf package name")
	buildCmd.Flags().StringToStringVarP(&protobuf.Options, "options", "t", nil, "protobuf options key value")
	buildCmd.MarkFlagRequired("input") // nolint
	buildCmd.MarkFlagRequired("out")   // nolint
	buildCmd.MarkFlagsRequiredTogether(
		"enabled",
		"dir",
		"package",
		"options",
	)

	buildProtoSubCmd.Flags().BoolVarP(&protobuf.Merge, "merge", "m", false, "protobuf merge in a file or not(default: false)")
	buildProtoSubCmd.Flags().StringVarP(&protobuf.Dir, "dir", "d", "", "protobuf out directory")
	buildProtoSubCmd.Flags().StringVarP(&protobuf.Package, "package", "p", "", "protobuf package name")
	buildProtoSubCmd.Flags().StringToStringVarP(&protobuf.Options, "options", "t", nil, "protobuf options key value")
	buildProtoSubCmd.MarkFlagsRequiredTogether(
		"dir",
		"package",
		"options",
	)

	buildCmd.AddCommand(buildProtoSubCmd)
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
		c.OutDir = outDir
		c.View.Protobuf = view.Protobuf{
			Enabled: protobuf.Enabled,
			Dir:     protobuf.Dir,
			Merge:   protobuf.Merge,
			Package: protobuf.Package,
			Options: protobuf.Options,
		}
		setupBase(c)

		protobufConfig := &c.View.Protobuf

		generateModelFile := func(filename string, buf *bytes.Buffer) error {
			content, err := os.ReadFile(filename)
			if err != nil {
				return err
			}
			vw := view.New(&driver.SQL{
				CreateTableSQL:   string(content),
				CustomDefineType: c.TypeDefine,
			}, c.View)
			list, err := vw.GetDbFile(utils.GetPkgName(c.OutDir))
			if err != nil {
				return err
			}
			for _, v := range list {
				path := c.OutDir + "/" + v.Filename + ".go"
				_ = utils.WriteFile(path, v.Build())

				cmd, _ := exec.Command("goimports", "-l", "-w", path).Output()
				_, _ = exec.Command("gofmt", "-l", "-w", path).Output()
				log.Info("👉 " + strings.TrimSuffix(string(cmd), "\n"))

				if protobufConfig.Enabled {
					if protobufConfig.Merge {
						content := v.BuildProtobufEnumBody()
						if len(content) > 0 {
							buf.Write(content)
						}
					} else {
						content := v.BuildProtobufEnum()
						if len(content) > 0 {
							protoFilename := protobufConfig.Dir + "/" + v.Filename + ".proto"
							_ = utils.WriteFile(protoFilename, content)
							log.Info("👆 " + protoFilename)
						}
					}
				}
			}
			return nil
		}

		buf := &bytes.Buffer{}
		for _, filename := range inputFile {
			err = generateModelFile(filename, buf)
			if err != nil {
				log.Warnf("🧐 generate file from SQL file(%s) failed !!!", filename)
			}
		}
		if protobufConfig.Enabled && protobufConfig.Merge && buf.Len() > 0 {
			filename := utils.GetPkgName(protobufConfig.Dir)
			protoFilename := protobufConfig.Dir + "/" + filename + ".proto"
			header := ast.BuildRawProtobufEnumHeader(protobufConfig.Package, protobufConfig.Options)
			_ = utils.WriteFile(protoFilename, append(header, buf.Bytes()...))
			log.Info("👆 " + protoFilename)
		}

		log.Info("😄 generate success !!!")
		return nil
	},
}

var buildProtoSubCmd = &cobra.Command{
	Use:     "proto",
	Short:   "Generate model from sql",
	Example: "ormat build proto",
	RunE: func(*cobra.Command, []string) error {
		c := config.Global
		err := c.Load()
		if err != nil {
			return err
		}
		c.OutDir = outDir
		c.View.Protobuf = view.Protobuf{
			Enabled: true,
			Dir:     protobuf.Dir,
			Merge:   protobuf.Merge,
			Package: protobuf.Package,
			Options: protobuf.Options,
		}
		setupBase(c)

		protobufConfig := &c.View.Protobuf

		generateModelFile := func(filename string, buf *bytes.Buffer) error {
			content, err := os.ReadFile(filename)
			if err != nil {
				return err
			}
			vw := view.New(&driver.SQL{
				CreateTableSQL:   string(content),
				CustomDefineType: c.TypeDefine,
			}, c.View)
			list, err := vw.GetDbFile(utils.GetPkgName(c.OutDir))
			if err != nil {
				return err
			}
			for _, v := range list {
				if protobufConfig.Merge {
					content := v.BuildProtobufEnumBody()
					if len(content) > 0 {
						buf.Write(content)
					}
				} else {
					content := v.BuildProtobufEnum()
					if len(content) > 0 {
						protoFilename := protobufConfig.Dir + "/" + v.Filename + ".proto"
						_ = utils.WriteFile(protoFilename, content)
						log.Info("👆 " + protoFilename)
					}
				}
			}
			return nil
		}

		buf := &bytes.Buffer{}
		for _, filename := range inputFile {
			err = generateModelFile(filename, buf)
			if err != nil {
				log.Warnf("🧐 generate file from SQL file(%s) failed !!!", filename)
			}
		}
		if protobufConfig.Enabled && protobufConfig.Merge && buf.Len() > 0 {
			filename := utils.GetPkgName(protobufConfig.Dir)
			protoFilename := protobufConfig.Dir + "/" + filename + ".proto"
			header := ast.BuildRawProtobufEnumHeader(protobufConfig.Package, protobufConfig.Options)
			_ = utils.WriteFile(protoFilename, append(header, buf.Bytes()...))
			log.Info("👆 " + protoFilename)
		}

		log.Info("😄 generate success !!!")
		return nil
	},
}
