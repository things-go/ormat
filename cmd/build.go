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
	InputFile     []string
	OutputDir     string
	Merge         bool
	MergeFilename string
	Package       string
	Options       map[string]string
	Suffix        string
	Template      string
	TypeDefine    map[string]string
	View          view.Config
}

type buildCmd struct {
	cmd *cobra.Command
	buildOpt
}

func newBuildCmd() *buildCmd {
	root := &buildCmd{
		buildOpt: buildOpt{
			View: view.Config{
				DbTag:            "gorm",
				Tags:             map[string]string{"json": view.TagSnakeCase},
				EnableLint:       false,
				DisableNull:      false,
				EnableInt:        false,
				EnableIntegerInt: false,
				EnableBoolInt:    false,
				IsNullToPoint:    true,
				IsForeignKey:     false,
				IsCommentTag:     true,
			},
		},
	}

	PreRunBuild := func(*cobra.Command, []string) error {
		setupBase2("prod")
		return nil
	}
	cmd := &cobra.Command{
		Use:     "build",
		Short:   "Generate model from sql",
		Example: "ormat build",
		PreRunE: PreRunBuild,
		RunE: func(*cobra.Command, []string) error {
			usedTemplateMapping, err := getModelTemplate(root.Template, root.Suffix)
			if err != nil {
				return err
			}
			genFile := &generateFile{
				Files:         parseSqlFromFile(&root.buildOpt),
				OutputDir:     root.OutputDir,
				Template:      usedTemplateMapping.Template,
				Merge:         root.Merge,
				MergeFilename: root.MergeFilename,
				Package:       root.Package,
				Options:       root.Options,
				Suffix:        usedTemplateMapping.Suffix,
				GenFunc:       genModelFile,
			}
			genFile.runGenModel()
			return nil
		},
	}

	cmdInfo := &cobra.Command{
		Use:     "info",
		Short:   "model info from sql",
		Example: "ormat build info",
		PreRunE: PreRunBuild,
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
			genFile.runGenModel()
			return nil
		},
	}

	cmdEnum := &cobra.Command{
		Use:     "enum",
		Short:   "Generate enum from sql",
		Example: "ormat build enum",
		PreRunE: PreRunBuild,
		RunE: func(*cobra.Command, []string) error {
			usedTemplate, err := getEnumTemplate(root.Template, root.Suffix)
			if err != nil {
				return err
			}
			genFile := &generateFile{
				Files:         parseSqlFromFile(&root.buildOpt),
				Template:      usedTemplate.Template,
				OutputDir:     root.OutputDir,
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

	cmdEnumInfo := &cobra.Command{
		Use:     "info",
		Short:   "enum info from sql",
		Example: "ormat build enum info",
		PreRunE: PreRunBuild,
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
			genFile.runGenEnum()
			return nil
		},
	}

	cmd.PersistentFlags().StringSliceVarP(&root.InputFile, "input", "i", nil, "input file")
	cmd.PersistentFlags().StringVarP(&root.OutputDir, "out", "o", "", "out directory")
	cmd.PersistentFlags().BoolVar(&root.Merge, "merge", false, "merge in a file or not")
	cmd.PersistentFlags().StringVar(&root.MergeFilename, "filename", "", "merge filename")
	cmd.PersistentFlags().StringVar(&root.Package, "package", "", "package name")
	cmd.PersistentFlags().StringToStringVar(&root.Options, "options", nil, "options key value")
	cmd.PersistentFlags().StringVar(&root.Suffix, "suffix", "", "filename suffix")
	cmd.PersistentFlags().StringVar(&root.Template, "template", "__in_go", "use custom template")
	cmd.MarkPersistentFlagRequired("input") // nolint
	cmd.MarkPersistentFlagRequired("out")   // nolint

	cmdEnum.AddCommand(
		cmdEnumInfo,
	)
	cmd.AddCommand(
		cmdInfo,
		cmdEnum,
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
	astFiles := make([]*ast.File, 0, 64)
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
