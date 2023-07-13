package command

import (
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
	DbConfig
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
			db, dbName, err := NewDB(root.DbConfig)
			if err != nil {
				return err
			}
			defer CloseDB(db)
			d, err := NewDriver(&DriverConfig{
				DB:         db,
				Dialect:    root.Dialect,
				DbName:     dbName,
				TableNames: nil,
			})
			if err != nil {
				return err
			}
			mixin, err := d.GetSchema()
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
					data := codegen.New([]*ens.Entity{entity}, codegenOption...).
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
	cmd.Flags().StringVar(&root.DbConfig.Dialect, "dialect", "mysql", "database dialect, one of [mysql,sqlite3]")
	cmd.Flags().StringVar(&root.DbConfig.DSN, "dsn", "", "database dsn(root:123456@tcp(127.0.0.1:3306)/test)")
	cmd.Flags().StringVar(&root.DbConfig.Options, "option", "", "database option(dsn?option)")
	cmd.Flags().StringVarP(&root.OutputDir, "out", "o", "./model/migration", "out directory")
	cmd.Flags().StringVar(&root.Filename, "filename", "migration", "filename")

	cmd.MarkFlagRequired("dsn")

	root.cmd = cmd
	return root
}
