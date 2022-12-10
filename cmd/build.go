package cmd

import (
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/things-go/log"

	"github.com/things-go/ormat/pkg/config"
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
	Suffix:        "",
	Template:      "",
}

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
	buildEnumCmd.PersistentFlags().StringVarP(&protobuf.Suffix, "suffix", "s", ".proto", "out filename suffix")
	buildEnumCmd.MarkFlagsRequiredTogether(
		"dir",
		"package",
		"options",
	)

	buildEnumCustomCmd.Flags().StringVar(&protobuf.Template, "template", "", "use custom template")
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
	PreRunE: PreRunBuildEnum,
	RunE: func(*cobra.Command, []string) error {
		c := config.Global

		astFiles := parseSqlFromFile(c, inputFile)
		genModelFile(astFiles, &c.View.Protobuf, c.OutDir)
		return nil
	},
}

var buildEnumCmd = &cobra.Command{
	Use:     "enum",
	Short:   "Generate enum from sql",
	Example: "ormat build enum",
	PreRunE: PreRunBuildEnum,
	RunE: func(*cobra.Command, []string) error {
		return runGenEnumFile(tpl.ProtobufEnumTpl)
	},
}

var buildEnumMappingCmd = &cobra.Command{
	Use:     "mapping",
	Short:   "Generate enum mapping from sql",
	Example: "ormat build enum mapping",
	PreRunE: PreRunBuildEnum,
	RunE: func(*cobra.Command, []string) error {
		return runGenEnumFile(tpl.ProtobufEnumMappingTpl)
	},
}

var buildEnumCustomCmd = &cobra.Command{
	Use:     "custom",
	Short:   "Generate enum custom with template from sql",
	Example: "ormat build enum custom",
	PreRunE: PreRunBuildEnum,
	RunE: func(*cobra.Command, []string) error {
		usedTemplate, err := parseTemplateFromFile(protobuf.Template)
		if err != nil {
			return err
		}
		return runGenEnumFile(usedTemplate)
	},
}

func PreRunBuildEnum(*cobra.Command, []string) error {
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
		Suffix:  GetFilenameSuffix(protobuf.Suffix),
	}
	setupBase(c)
	return nil
}

func runGenEnumFile(usedTemplate *template.Template) error {
	c := config.Global

	astFiles := parseSqlFromFile(c, inputFile)
	genEnumFile(astFiles, &c.View.Protobuf, usedTemplate)
	return nil
}

func parseSqlFromFile(c *config.Config, inputFiles []string) []*ast.File {
	innerParseFromFile := func(filename string) ([]*ast.File, error) {
		content, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		return view.New(
			&driver.SQL{
				CreateTableSQL:   string(content),
				CustomDefineType: c.TypeDefine,
			},
			c.View,
		).GetDbFile(utils.GetPkgName(c.OutDir))
	}
	astFiles := make([]*ast.File, 0, 64)
	for _, filename := range inputFiles {
		astFile, err := innerParseFromFile(filename)
		if err != nil {
			log.Warnf("🧐 parse from SQL file(%s) failed !!!", filename)
			continue
		}
		astFiles = append(astFiles, astFile...)
	}
	return astFiles
}
