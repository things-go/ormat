package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/things-go/log"

	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/view"
	"github.com/things-go/ormat/view/ast"
	"github.com/things-go/ormat/view/driver"
)

type buildOpt struct {
	InputFile  []string
	OutputDir  string
	TypeDefine map[string]string
	View       view.Config

	Merge         bool
	MergeFilename string
	Package       string
	Options       map[string]string
	Suffix        string
	Template      string
}

type buildCmd struct {
	cmd *cobra.Command
	buildOpt
}

func newBuildCmd() *buildCmd {
	root := &buildCmd{}

	cmd := &cobra.Command{
		Use:     "build",
		Short:   "Generate model from sql",
		Example: "ormat build",
		RunE: func(*cobra.Command, []string) error {
			usedTemplate, err := getModelTemplate(root.Template, root.Suffix)
			if err != nil {
				return err
			}
			genFile := &generateFile{
				Files:         parseSqlFromFile(&root.buildOpt),
				OutputDir:     root.OutputDir,
				Template:      usedTemplate.Template,
				Merge:         root.Merge,
				MergeFilename: root.MergeFilename,
				Package:       root.Package,
				Options:       root.Options,
				Suffix:        usedTemplate.Suffix,
				GenFunc:       genModelFile,
			}
			genFile.runGen()
			return nil
		},
	}

	cmdInfo := &cobra.Command{
		Use:     "info",
		Short:   "model info from sql",
		Example: "ormat build info",
		RunE: func(*cobra.Command, []string) error {
			genFile := &generateFile{
				Files:         parseSqlFromFile(&root.buildOpt),
				Template:      nil,
				OutputDir:     root.OutputDir,
				Merge:         true,
				MergeFilename: root.MergeFilename,
				Package:       root.Package,
				Options:       root.Options,
				Suffix:        root.Suffix,
				GenFunc:       showInformation,
			}
			genFile.runGen()
			return nil
		},
	}

	cmd.PersistentFlags().StringSliceVarP(&root.InputFile, "input", "i", nil, "input file")
	cmd.PersistentFlags().StringVarP(&root.OutputDir, "out", "o", "./model", "out directory")

	cmd.Flags().StringToStringVarP(&root.TypeDefine, "define", "D", nil, "custom type define")
	view.InitFlagSetForConfig(cmd.Flags(), &root.View)

	cmd.PersistentFlags().BoolVar(&root.Merge, "merge", false, "merge in a file or not")
	cmd.PersistentFlags().StringVar(&root.MergeFilename, "filename", "", "merge filename")
	cmd.PersistentFlags().StringVar(&root.Package, "package", "", "package name")
	cmd.PersistentFlags().StringToStringVar(&root.Options, "options", nil, "options key value")
	cmd.PersistentFlags().StringVar(&root.Suffix, "suffix", "", "filename suffix")
	cmd.PersistentFlags().StringVar(&root.Template, "template", "__in_go", "use custom template")

	cmd.MarkPersistentFlagRequired("input") // nolint
	cmd.AddCommand(
		cmdInfo,
	)
	root.cmd = cmd
	return root
}

func parseSqlFromFile(c *buildOpt) []*ast.File {
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
		).GetDbFile(utils.GetPkgName(c.OutputDir))
	}
	astFiles := make([]*ast.File, 0, 256)
	for _, filename := range c.InputFile {
		astFile, err := innerParseFromFile(filename)
		if err != nil {
			log.Warnf("🧐 parse from SQL file(%s) failed !!!", filename)
			continue
		}
		astFiles = append(astFiles, astFile...)
	}
	return astFiles
}
