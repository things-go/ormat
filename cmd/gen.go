package cmd

import (
	"github.com/spf13/cobra"

	"github.com/things-go/ormat/pkg/database"
	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/runtime"
	"github.com/things-go/ormat/view"
	"github.com/things-go/ormat/view/ast"
)

type genOpt struct {
	OutputDir string
	runtime.Database
	TableNames []string
	TypeDefine map[string]string
	View       view.Config

	Merge         bool
	MergeFilename string
	Package       string
	Options       map[string]string
	Suffix        string
	Template      string
}

type genCmd struct {
	cmd *cobra.Command
	genOpt
}

func newGenCmd() *genCmd {
	root := &genCmd{}

	getAstFiles := func() ([]*ast.File, error) {
		err := root.Database.Parse()
		if err != nil {
			return nil, err
		}
		db, err := runtime.NewDb(&root.Database)
		if err != nil {
			return nil, err
		}
		defer database.Close(db)

		vw, err := NewFromDatabase(&DbConfig{
			DB:               db,
			Dialect:          root.Database.Dialect,
			DbName:           root.DbName(),
			TableNames:       root.TableNames,
			CustomDefineType: root.TypeDefine,
			Config:           root.View,
		})
		if err != nil {
			return nil, err
		}
		return vw.GetDbFile(utils.GetPkgName(root.OutputDir))
	}

	cmd := &cobra.Command{
		Use:     "gen",
		Short:   "Generate model from database",
		Example: "ormat gen",
		RunE: func(*cobra.Command, []string) error {
			usedTemplate, err := getModelTemplate(root.Template, root.Suffix)
			if err != nil {
				return err
			}
			files, err := getAstFiles()
			if err != nil {
				return err
			}
			genFile := &generateFile{
				Files:         files,
				OutputDir:     root.OutputDir,
				Template:      usedTemplate.Template,
				Merge:         root.Merge,
				MergeFilename: root.MergeFilename,
				Package:       root.Package,
				Options:       root.Options,
				Suffix:        usedTemplate.Suffix,
				GenFunc:       genModelFile,
			}
			genFile.runGenModel()
			return nil
		},
	}
	cmdInfo := &cobra.Command{
		Use:     "info",
		Short:   "model info from database",
		Example: "ormat gen info",
		RunE: func(*cobra.Command, []string) error {
			files, err := getAstFiles()
			if err != nil {
				return err
			}
			genFile := &generateFile{
				Files:         files,
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
		Short:   "Generate enum from database",
		Example: "ormat gen enum",
		RunE: func(*cobra.Command, []string) error {
			usedTemplate, err := getEnumTemplate(root.Template, root.Suffix)
			if err != nil {
				return err
			}
			files, err := getAstFiles()
			if err != nil {
				return err
			}
			genFile := &generateFile{
				Files:         files,
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
		Short:   "enum info from database",
		Example: "ormat gen enum info",
		RunE: func(*cobra.Command, []string) error {
			files, err := getAstFiles()
			if err != nil {
				return err
			}
			genFile := &generateFile{
				Files:         files,
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

	cmd.PersistentFlags().StringVar(&root.Database.Dialect, "dialect", "mysql", "database dialect, one of [mysql,sqlite3]")
	cmd.PersistentFlags().StringVar(&root.Database.DSN, "dsn", "", "database dsn(root:123456@tcp(127.0.0.1:3306)/test)")
	cmd.PersistentFlags().StringVar(&root.Database.Options, "option", "", "database option(dsn?option)")
	cmd.PersistentFlags().StringVarP(&root.OutputDir, "out", "o", "./model", "out directory")
	cmd.PersistentFlags().StringSliceVarP(&root.TableNames, "table", "t", nil, "only out custom table")

	cmd.Flags().StringToStringVarP(&root.TypeDefine, "define", "D", nil, "custom type define")
	view.InitFlagSetForConfig(cmd.Flags(), &root.View)

	cmd.PersistentFlags().BoolVar(&root.Merge, "merge", false, "merge in a file or not")
	cmd.PersistentFlags().StringVar(&root.MergeFilename, "filename", "", "merge filename")
	cmd.PersistentFlags().StringVar(&root.Package, "package", "", "package name")
	cmd.PersistentFlags().StringToStringVar(&root.Options, "options", nil, "options key value")
	cmd.PersistentFlags().StringVar(&root.Suffix, "suffix", "", "filename suffix")
	cmd.PersistentFlags().StringVar(&root.Template, "template", "__in_go", "use custom template")

	cmd.MarkPersistentFlagRequired("dsn") // nolint

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
