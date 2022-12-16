package cmd

import (
	"github.com/spf13/cobra"
	"github.com/things-go/log"

	"github.com/things-go/ormat/pkg/database"
	"github.com/things-go/ormat/pkg/tpl"
	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/runtime"
	"github.com/things-go/ormat/view"
)

type sqlOpt struct {
	OutputDir string
	Filename  string
	runtime.Database
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
			err := root.Database.Parse()
			if err != nil {
				return err
			}
			db, err := runtime.NewDb(&root.Database)
			if err != nil {
				return err
			}
			defer database.Close(db)
			vw, err := NewFromDatabase(&DbConfig{
				DB:               db,
				Dialect:          root.Database.Dialect,
				DbName:           root.DbName(),
				TableNames:       nil,
				CustomDefineType: nil,
				Config:           view.Config{},
			})
			if err != nil {
				return err
			}
			sqlFile, err := vw.GetSqlFile()
			if err != nil {
				return err
			}
			filename := intoFilename(root.OutputDir, root.Filename, ".sql")
			err = utils.WriteFileWithTemplate(filename, tpl.SqlDDL, sqlFile)
			if err != nil {
				return err
			}
			log.Info("👉 " + filename)
			return nil
		},
	}
	cmd.Flags().StringVar(&root.Database.Dialect, "dialect", "mysql", "database dialect, one of [mysql,sqlite3]")
	cmd.Flags().StringVar(&root.Database.DSN, "dsn", "", "database dsn(root:123456@tcp(127.0.0.1:3306)/test)")
	cmd.Flags().StringVar(&root.Database.Options, "option", "", "database option(dsn?option)")
	cmd.Flags().StringVarP(&root.OutputDir, "out", "o", "./model", "out directory")
	cmd.Flags().StringVar(&root.Filename, "filename", "migration", "filename")

	cmd.MarkFlagRequired("dsn")

	root.cmd = cmd
	return root
}
