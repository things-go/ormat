package command

import (
	"context"

	"ariga.io/atlas/sql/schema"
	"github.com/spf13/cobra"
	"github.com/things-go/ens"
	"github.com/things-go/ens/codegen"
	"github.com/things-go/ens/utils"
	"golang.org/x/exp/slog"
)

type sqlOpt struct {
	OutputDir string
	Merge     bool
	Filename  string
	URL       string
	Tables    []string
	Exclude   []string
}

type sqlCmd struct {
	cmd *cobra.Command
	sqlOpt
}

func newSqlCmd() *sqlCmd {
	root := &sqlCmd{}
	cmd := &cobra.Command{
		Use:     "sql",
		Short:   "Generate sql file",
		Example: "ormat sql",
		RunE: func(*cobra.Command, []string) error {
			d, err := NewDriver(root.URL)
			if err != nil {
				return err
			}
			mixin, err := d.InspectSchema(context.Background(), &schema.InspectOptions{
				Mode:    schema.InspectTables,
				Tables:  root.Tables,
				Exclude: root.Exclude,
			})
			if err != nil {
				return err
			}
			sc := mixin.Build(nil)
			codegenOption := []codegen.Option{
				codegen.WithByName("ormat"),
				codegen.WithVersion(version),
				codegen.WithPackageName(utils.GetPkgName(root.OutputDir)),
			}
			if root.Merge {
				data := codegen.New(sc.Entities, codegenOption...).
					GenDDL().
					Bytes()
				filename := joinFilename(root.OutputDir, root.Filename, ".sql")
				err = WriteFile(filename, data)
				if err != nil {
					return err
				}
				slog.Info("👉 " + filename)
			} else {
				for _, entity := range sc.Entities {
					data := codegen.New([]*ens.EntityDescriptor{entity}, codegenOption...).
						GenDDL().
						Bytes()
					filename := joinFilename(root.OutputDir, entity.Name, ".sql")
					err = WriteFile(filename, data)
					if err != nil {
						return err
					}
					slog.Info("👉 " + filename)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root.URL, "url", "", "mysql://root:123456@127.0.0.1:3306/test)")
	cmd.PersistentFlags().StringSliceVarP(&root.Tables, "table", "t", nil, "only out custom table")
	cmd.PersistentFlags().StringSliceVar(&root.Exclude, "exclude", nil, "exclude table pattern")
	cmd.Flags().StringVarP(&root.OutputDir, "out", "o", "./model/migration", "out directory")
	cmd.Flags().StringVar(&root.Filename, "filename", "migration", "filename when merge enabled")
	cmd.Flags().BoolVar(&root.Merge, "merge", false, "merge in a file")

	cmd.MarkFlagRequired("url")

	root.cmd = cmd
	return root
}
