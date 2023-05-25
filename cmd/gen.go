package cmd

import (
	"github.com/spf13/cobra"

	"github.com/things-go/ormat/pkg/database"
	"github.com/things-go/ormat/pkg/tpl"
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
	Suffix        string
	Template      string
	HasAssist     bool // 是否提供辅助工具集
	HasEntity     bool // 提供一般性查询
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
				Files:          files,
				OutputDir:      root.OutputDir,
				Template:       usedTemplate.Template,
				Merge:          root.Merge,
				MergeFilename:  root.MergeFilename,
				Package:        root.View.Package,
				Options:        root.View.Options,
				Suffix:         usedTemplate.Suffix,
				GenFunc:        genModelFile,
				HasAssist:      root.HasAssist,
				AssistTemplate: tpl.Assist,
				AssistGenFunc:  genModelFile,
			}
			genFile.runGen()
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
				Package:       root.View.Package,
				Options:       root.View.Options,
				Suffix:        root.Suffix,
				GenFunc:       showInformation,
			}
			genFile.runGen()
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
	cmd.PersistentFlags().StringVar(&root.Suffix, "suffix", "", "filename suffix")
	cmd.PersistentFlags().StringVar(&root.Template, "template", "__in_go", "use custom template")
	cmd.PersistentFlags().BoolVar(&root.HasAssist, "hasAssist", false, "是否提供辅助工具集")
	cmd.PersistentFlags().BoolVar(&root.HasEntity, "hasEntity", false, "是否提供查询工具集")

	cmd.MarkPersistentFlagRequired("dsn") // nolint
	cmd.AddCommand(
		cmdInfo,
	)
	root.cmd = cmd
	return root
}
