package cmd

import (
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/things-go/log"

	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/pkg/consts"
	"github.com/things-go/ormat/pkg/tpl"
	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/view"
	"github.com/things-go/ormat/view/ast"
	"github.com/things-go/ormat/view/driver"
)

var inputFile []string
var outDir string
var protobuf = view.Protobuf{
	Enabled:       false,
	Merge:         false,
	MergeFilename: "",
	Dir:           "",
	Package:       "",
	Options:       nil,
}
var filenameSuffix string
var customTemplate string

func init() {
	buildCmd.PersistentFlags().StringSliceVarP(&inputFile, "input", "i", nil, "input file")

	buildCmd.Flags().StringVarP(&outDir, "out", "o", "", "model out directory")
	buildCmd.Flags().BoolVarP(&protobuf.Enabled, "enabled", "e", false, "protobuf enabled or not(default: false)")
	buildCmd.Flags().BoolVarP(&protobuf.Merge, "merge", "m", false, "protobuf merge in a file or not(default: false)")
	buildCmd.Flags().StringVarP(&protobuf.MergeFilename, "filename", "f", "", "merge filename")
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

	buildEnumCmd.PersistentFlags().BoolVarP(&protobuf.Merge, "merge", "m", false, "protobuf merge in a file or not(default: false)")
	buildEnumCmd.PersistentFlags().StringVarP(&protobuf.Dir, "dir", "d", "", "protobuf out directory")
	buildEnumCmd.PersistentFlags().StringVarP(&protobuf.Package, "package", "p", "", "protobuf package name")
	buildEnumCmd.PersistentFlags().StringToStringVarP(&protobuf.Options, "options", "t", nil, "protobuf options key value")
	buildEnumCmd.PersistentFlags().StringVarP(&filenameSuffix, "suffix", "s", ".proto", "out filename suffix")
	buildEnumCmd.MarkFlagsRequiredTogether(
		"dir",
		"package",
		"options",
	)

	buildEnumCustomCmd.Flags().StringVar(&customTemplate, "template", "", "use custom template")
	buildEnumCustomCmd.MarkFlagRequired("template")

	buildEnumCmd.AddCommand(
		buildEnumMappingCmd,
		buildEnumCustomCmd,
	)
	buildCmd.AddCommand(buildEnumCmd)
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
		mergeProtoEnumFile := ast.ProtobufEnumFile{
			Version: consts.Version,
			Package: protobufConfig.Package,
			Options: protobufConfig.Options,
			Enums:   make([]*ast.ProtobufEnum, 0, 64),
		}

		generateModelFile := func(filename string) error {
			content, err := os.ReadFile(filename)
			if err != nil {
				return err
			}
			vw := view.New(
				&driver.SQL{
					CreateTableSQL:   string(content),
					CustomDefineType: c.TypeDefine,
				},
				c.View,
			)
			list, err := vw.GetDbFile(utils.GetPkgName(c.OutDir))
			if err != nil {
				return err
			}
			for _, v := range list {
				path := c.OutDir + "/" + v.Filename + ".go"
				_ = utils.WriteFileWithTemplate(path, tpl.ModelTpl, v)

				cmd, _ := exec.Command("goimports", "-l", "-w", path).Output()
				_, _ = exec.Command("gofmt", "-l", "-w", path).Output()
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
			return nil
		}

		for _, filename := range inputFile {
			err = generateModelFile(filename)
			if err != nil {
				log.Warnf("🧐 generate file from SQL file(%s) failed !!!", filename)
			}
		}
		if protobufConfig.Enabled && protobufConfig.Merge && len(mergeProtoEnumFile.Enums) > 0 {
			enumFilename := intoFilename(protobufConfig.Dir, protobufConfig.GetMergeFilename(), ".proto")
			_ = utils.WriteFileWithTemplate(enumFilename, tpl.ProtobufEnumTpl, mergeProtoEnumFile)
			log.Info("👆 " + enumFilename)
		}

		log.Info("😄 generate success !!!")
		return nil
	},
}

var buildEnumCmd = &cobra.Command{
	Use:     "enum",
	Short:   "Generate enum from sql",
	Example: "ormat build enum",
	RunE: func(*cobra.Command, []string) error {
		return generateEnumFile(tpl.ProtobufEnumTpl)
	},
}

var buildEnumMappingCmd = &cobra.Command{
	Use:     "mapping",
	Short:   "Generate enum mapping from sql",
	Example: "ormat build enum mapping",
	RunE: func(*cobra.Command, []string) error {
		return generateEnumFile(tpl.ProtobufEnumMappingTpl)
	},
}

var buildEnumCustomCmd = &cobra.Command{
	Use:     "custom",
	Short:   "Generate enum custom with template from sql",
	Example: "ormat build enum custom",
	RunE: func(*cobra.Command, []string) error {
		usedTemplate, err := parseTemplateFromFile(customTemplate)
		if err != nil {
			return err
		}
		return generateEnumFile(usedTemplate)
	},
}

func generateEnumFile(usedTemplate *template.Template) error {
	c := config.Global
	err := c.Load()
	if err != nil {
		return err
	}
	c.View.Protobuf = view.Protobuf{
		Enabled: true,
		Dir:     protobuf.Dir,
		Merge:   protobuf.Merge,
		Package: protobuf.Package,
		Options: protobuf.Options,
	}
	setupBase(c)
	filenameSuffix = GetFilenameSuffix(filenameSuffix)
	protobufConfig := &c.View.Protobuf

	mergeProtoEnumFile := ast.ProtobufEnumFile{
		Version: consts.Version,
		Package: protobufConfig.Package,
		Options: protobufConfig.Options,
		Enums:   nil,
	}

	genEnumFile := func(filename string) error {
		content, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		vw := view.New(
			&driver.SQL{
				CreateTableSQL:   string(content),
				CustomDefineType: c.TypeDefine,
			},
			c.View,
		)
		list, err := vw.GetDbFile(utils.GetPkgName(c.OutDir))
		if err != nil {
			return err
		}
		for _, v := range list {
			if enums := v.GetProtobufEnums(); len(enums) > 0 {
				if protobufConfig.Merge {
					mergeProtoEnumFile.Enums = append(mergeProtoEnumFile.Enums, enums...)
				} else {
					protoEnumFile := ast.ProtobufEnumFile{
						Version: consts.Version,
						Package: protobufConfig.Package,
						Options: protobufConfig.Options,
						Enums:   enums,
					}
					enumFilename := intoFilename(protobufConfig.Dir, v.Filename, filenameSuffix)
					_ = utils.WriteFileWithTemplate(enumFilename, usedTemplate, protoEnumFile)
					log.Info("👆 " + enumFilename)
				}
			}
		}
		return nil
	}

	for _, filename := range inputFile {
		err = genEnumFile(filename)
		if err != nil {
			log.Warnf("🧐 generate file from SQL file(%s) failed !!!", filename)
		}
	}

	if protobufConfig.Merge && len(mergeProtoEnumFile.Enums) > 0 {
		mergeFilename := protobufConfig.GetMergeFilename()
		enumFilename := intoFilename(protobufConfig.Dir, mergeFilename, filenameSuffix)
		_ = utils.WriteFileWithTemplate(enumFilename, usedTemplate, mergeProtoEnumFile)
		log.Info("👆 " + enumFilename)
	}

	log.Info("😄 generate success !!!")
	return nil
}
