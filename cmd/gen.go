package cmd

import (
	"github.com/spf13/cobra"

	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/pkg/database"
	"github.com/things-go/ormat/pkg/tpl"
	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/runtime"
	"github.com/things-go/ormat/view"
	"github.com/things-go/ormat/view/ast"
)

type genCmd struct {
	cmd           *cobra.Command
	inputFile     []string
	outputDir     string
	Enabled       bool
	Merge         bool
	MergeFilename string
	Package       string
	Options       map[string]string
	Suffix        string
	Template      string
}

func newGenCmd() *genCmd {
	root := &genCmd{}
	cmd := &cobra.Command{
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
				Template:  tpl.Model,
				GenFunc:   genModelFile,
				Suffix:    ".go",
			}
			genFile.runGenModel()
			return nil
		},
	}

	cmdEnum := &cobra.Command{
		Use:     "enum",
		Short:   "Generate enum from database",
		Example: "ormat gen enum",
		RunE: func(*cobra.Command, []string) error {
			usedTemplate, err := getEnumTemplate(root.Template, root.Suffix)
			if err != nil {
				return err
			}
			files, err := getAstFileFromDatabase(root.outputDir)
			if err != nil {
				return err
			}
			c := config.Global
			genFile := &generateFile{
				Files:         files,
				Template:      usedTemplate.Template,
				OutputDir:     c.OutDir,
				Merge:         root.Merge,
				MergeFilename: root.MergeFilename,
				Package:       root.Package,
				Options:       root.Options,
				Suffix:        usedTemplate.Suffix,
				GenFunc:       genEnumFile,
			}
			genFile.runGenEnum()
			return nil
		},
	}

	cmdEnum.PersistentFlags().StringVarP(&root.outputDir, "out", "o", "", "out directory")
	cmdEnum.PersistentFlags().BoolVarP(&root.Merge, "merge", "m", false, "merge in a file or not")
	cmdEnum.PersistentFlags().StringVarP(&root.MergeFilename, "filename", "f", "", "merge filename")
	cmdEnum.PersistentFlags().StringVarP(&root.Package, "package", "p", "", "package name")
	cmdEnum.PersistentFlags().StringToStringVarP(&root.Options, "options", "t", nil, "options key value")
	cmdEnum.PersistentFlags().StringVarP(&root.Suffix, "suffix", "s", "", "filename suffix")
	cmdEnum.PersistentFlags().StringVar(&root.Template, "template", "__in_go", "use custom template")

	cmdEnum.MarkPersistentFlagRequired("dsn") // nolint
	cmdEnum.MarkPersistentFlagRequired("out") // nolint

	// cmdEnum.AddCommand(
	// 	genEnumInfoCmd,
	// )
	cmd.AddCommand(
		// genInfoCmd,
		cmdEnum,
	)

	root.cmd = cmd
	return root
}
func getAstFileFromDatabase(outputDir string) ([]*ast.File, error) {
	c := config.Global
	err := c.Load()
	if err != nil {
		return nil, err
	}
	c.OutDir = outputDir

	setupBase(c)

	rt, err := runtime.NewRuntime(c)
	if err != nil {
		return nil, err
	}
	defer database.Close(rt.DB)

	return view.New(GetViewModel(rt), c.View).
		GetDbFile(utils.GetPkgName(c.OutDir))
}
