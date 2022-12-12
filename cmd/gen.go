package cmd

import (
	"text/template"

	"github.com/spf13/cobra"

	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/pkg/database"
	"github.com/things-go/ormat/pkg/tpl"
	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/runtime"
	"github.com/things-go/ormat/view"
	"github.com/things-go/ormat/view/ast"
)

func init() {
	genEnumCmd.PersistentFlags().StringVarP(&outputDir, "out", "o", "", "out directory")
	genEnumCmd.PersistentFlags().BoolVarP(&protobuf.Merge, "merge", "m", false, "merge in a file or not")
	genEnumCmd.PersistentFlags().StringVarP(&protobuf.MergeFilename, "filename", "f", "", "merge filename")
	genEnumCmd.PersistentFlags().StringVarP(&protobuf.Package, "package", "p", "", "protobuf package name")
	genEnumCmd.PersistentFlags().StringToStringVarP(&protobuf.Options, "options", "t", nil, "protobuf options key value")
	genEnumCmd.PersistentFlags().StringVarP(&protobuf.Suffix, "suffix", "s", ".proto", "out filename suffix")
	genEnumCmd.MarkPersistentFlagRequired("dsn") // nolint
	genEnumCmd.MarkPersistentFlagRequired("out") // nolint

	genEnumCustomCmd.Flags().StringVar(&protobuf.Template, "template", "", "use custom template")
	genEnumCustomCmd.MarkFlagRequired("template")

	genEnumCmd.AddCommand(
		genEnumMappingCmd,
		genEnumCustomCmd,
		// genEnumInfoCmd,
	)
	genCmd.AddCommand(
		// genInfoCmd,
		genEnumCmd,
	)
}

var genCmd = &cobra.Command{
	Use:     "gen",
	Short:   "Generate model from database",
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

		astFiles, err := view.New(GetViewModel(rt), c.View).
			GetDbFile(utils.GetPkgName(c.OutDir))
		if err != nil {
			return err
		}
		genFile := &generateFile{
			Files:     astFiles,
			OutputDir: c.OutDir,
			Template:  tpl.ModelTpl,
			GenFunc:   genModelFile,
		}
		genFile.runGenModel()
		return nil
	},
}

var genEnumCmd = &cobra.Command{
	Use:     "enum",
	Short:   "Generate enum from database",
	Example: "ormat gen enum",
	RunE: func(*cobra.Command, []string) error {
		return runGenEnumFile(tpl.ProtobufEnumTpl)
	},
}

var genEnumMappingCmd = &cobra.Command{
	Use:     "mapping",
	Short:   "Generate enum mapping from database",
	Example: "ormat gen enum mapping",
	RunE: func(*cobra.Command, []string) error {
		return runGenEnumFile(tpl.ProtobufEnumMappingTpl)
	},
}

var genEnumCustomCmd = &cobra.Command{
	Use:     "custom",
	Short:   "Generate enum custom with template from database",
	Example: "ormat gen enum custom",
	RunE: func(*cobra.Command, []string) error {
		usedTemplate, err := parseTemplateFromFile(protobuf.Template)
		if err != nil {
			return err
		}
		return runGenEnumFile(usedTemplate)
	},
}

func runGenEnumFile(usedTemplate *template.Template) error {
	files, err := getAstFileFromDatabase()
	if err != nil {
		return err
	}
	return runGenEnum(files, usedTemplate)
}

func getAstFileFromDatabase() ([]*ast.File, error) {
	c := config.Global
	err := c.Load()
	if err != nil {
		return nil, err
	}
	c.OutDir = outputDir

	setupBase(c)
	protobuf.Suffix = intoFilenameSuffix(protobuf.Suffix, ".proto")

	rt, err := runtime.NewRuntime(c)
	if err != nil {
		return nil, err
	}
	defer database.Close(rt.DB)

	return view.New(GetViewModel(rt), c.View).
		GetDbFile(utils.GetPkgName(c.OutDir))
}
