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
var outputDir string
var protobuf = view.Protobuf{
	Enabled:       true,
	Merge:         false,
	MergeFilename: "",
	Package:       "",
	Options:       nil,
	Suffix:        "",
	Template:      "",
}

func init() {
	buildCmd.PersistentFlags().StringSliceVarP(&inputFile, "input", "i", nil, "input file")
	buildCmd.PersistentFlags().StringVarP(&outputDir, "out", "o", "", "out directory")
	buildCmd.MarkPersistentFlagRequired("input") // nolint
	buildCmd.MarkPersistentFlagRequired("out")   // nolint

	buildEnumCmd.PersistentFlags().BoolVarP(&protobuf.Merge, "merge", "m", false, "merge in a file or not")
	buildEnumCmd.PersistentFlags().StringVarP(&protobuf.MergeFilename, "filename", "f", "", "merge filename")
	buildEnumCmd.PersistentFlags().StringVarP(&protobuf.Package, "package", "p", "", "protobuf package name")
	buildEnumCmd.PersistentFlags().StringToStringVarP(&protobuf.Options, "options", "t", nil, "protobuf options key value")
	buildEnumCmd.PersistentFlags().StringVarP(&protobuf.Suffix, "suffix", "s", ".proto", "out filename suffix")

	buildEnumCustomCmd.Flags().StringVar(&protobuf.Template, "template", "", "use custom template")
	buildEnumCustomCmd.MarkFlagRequired("template")

	buildEnumCmd.AddCommand(
		buildEnumMappingCmd,
		buildEnumCustomCmd,
		buildEnumInfoCmd,
	)
	buildCmd.AddCommand(
		buildInfoCmd,
		buildEnumCmd,
	)
}

var buildCmd = &cobra.Command{
	Use:     "build",
	Short:   "Generate model from sql",
	Example: "ormat build",
	PreRunE: PreRunBuild,
	RunE: func(*cobra.Command, []string) error {
		c := config.Global

		genFile := &generateFile{
			Files:     parseSqlFromFile(c, inputFile),
			OutputDir: c.OutDir,
			Template:  tpl.ModelTpl,
			GenFunc:   genModelFile,
		}
		genFile.runGenModel()
		return nil
	},
}

var buildInfoCmd = &cobra.Command{
	Use:     "info",
	Short:   "model info from sql",
	Example: "ormat build info",
	PreRunE: PreRunBuild,
	RunE: func(*cobra.Command, []string) error {
		c := config.Global
		genFile := &generateFile{
			Files:         parseSqlFromFile(c, inputFile),
			Template:      nil,
			OutputDir:     c.OutDir,
			Merge:         protobuf.Merge,
			MergeFilename: protobuf.MergeFilename,
			Package:       protobuf.Package,
			Options:       protobuf.Options,
			Suffix:        protobuf.Suffix,
			GenFunc:       showInfo,
		}
		genFile.runGenModel()
		return nil
	},
}

var buildEnumCmd = &cobra.Command{
	Use:     "enum",
	Short:   "Generate enum from sql",
	Example: "ormat build enum",
	PreRunE: PreRunBuild,
	RunE: func(*cobra.Command, []string) error {
		return runBuildEnumFile(tpl.ProtobufEnumTpl)
	},
}

var buildEnumMappingCmd = &cobra.Command{
	Use:     "mapping",
	Short:   "Generate enum mapping from sql",
	Example: "ormat build enum mapping",
	PreRunE: PreRunBuild,
	RunE: func(*cobra.Command, []string) error {
		return runBuildEnumFile(tpl.ProtobufEnumMappingTpl)
	},
}

var buildEnumCustomCmd = &cobra.Command{
	Use:     "custom",
	Short:   "Generate enum custom with template from sql",
	Example: "ormat build enum custom",
	PreRunE: PreRunBuild,
	RunE: func(*cobra.Command, []string) error {
		usedTemplate, err := parseTemplateFromFile(protobuf.Template)
		if err != nil {
			return err
		}
		return runBuildEnumFile(usedTemplate)
	},
}

var buildEnumInfoCmd = &cobra.Command{
	Use:     "info",
	Short:   "enum info from sql",
	Example: "ormat build enum info",
	PreRunE: PreRunBuild,
	RunE: func(*cobra.Command, []string) error {
		c := config.Global
		genFile := &generateFile{
			Files:         parseSqlFromFile(c, inputFile),
			Template:      nil,
			OutputDir:     c.OutDir,
			Merge:         protobuf.Merge,
			MergeFilename: protobuf.MergeFilename,
			Package:       protobuf.Package,
			Options:       protobuf.Options,
			Suffix:        protobuf.Suffix,
			GenFunc:       showInfo,
		}
		genFile.runGenEnum()
		return nil
	},
}

func PreRunBuild(*cobra.Command, []string) error {
	c := config.Global
	err := c.Load()
	if err != nil {
		return err
	}
	c.OutDir = outputDir
	setupBase(c)

	protobuf.Suffix = intoFilenameSuffix(protobuf.Suffix, ".proto")
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

func runBuildEnumFile(usedTemplate *template.Template) error {
	c := config.Global
	files := parseSqlFromFile(c, inputFile)
	return runGenEnum(files, usedTemplate)
}

func runGenEnum(file []*ast.File, usedTemplate *template.Template) error {
	c := config.Global
	genFile := &generateFile{
		Files:         file,
		Template:      usedTemplate,
		OutputDir:     c.OutDir,
		Merge:         protobuf.Merge,
		MergeFilename: protobuf.MergeFilename,
		Package:       protobuf.Package,
		Options:       protobuf.Options,
		Suffix:        protobuf.Suffix,
		GenFunc:       genEnumFile,
	}
	genFile.runGenEnum()
	return nil
}
