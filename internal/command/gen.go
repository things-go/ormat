package command

import (
	"github.com/spf13/cobra"
	"github.com/things-go/ens"
)

type genOpt struct {
	DbConfig
	TableNames []string

	genFileOpt
}

type genCmd struct {
	cmd *cobra.Command
	genOpt
}

func newGenCmd() *genCmd {
	root := &genCmd{}

	getSchema := func() (ens.Schemaer, error) {
		db, dbName, err := NewDB(root.DbConfig)
		if err != nil {
			return nil, err
		}
		defer CloseDB(db)

		d, err := NewDriver(&DriverConfig{
			DB:         db,
			Dialect:    root.Dialect,
			DbName:     dbName,
			TableNames: root.TableNames,
		})
		if err != nil {
			return nil, err
		}
		return d.GetSchema()
	}

	cmd := &cobra.Command{
		Use:     "gen",
		Short:   "Generate model from database",
		Example: "ormat gen",
		RunE: func(*cobra.Command, []string) error {
			sc, err := getSchema()
			if err != nil {
				return err
			}
			return root.genFileOpt.GenModel(sc)
		},
	}

	cmdAssist := &cobra.Command{
		Use:     "assist",
		Short:   "model assist from database",
		Example: "ormat gen assist",
		RunE: func(*cobra.Command, []string) error {
			sc, err := getSchema()
			if err != nil {
				return err
			}
			return root.genFileOpt.GenAssist(sc)
		},
	}

	cmdMapper := &cobra.Command{
		Use:     "mapper",
		Short:   "model mapper from database",
		Example: "ormat gen mapper",
		RunE: func(*cobra.Command, []string) error {
			sc, err := getSchema()
			if err != nil {
				return err
			}
			return root.genFileOpt.GenMapper(sc)
		},
	}

	cmd.PersistentFlags().StringVar(&root.DbConfig.Dialect, "dialect", "mysql", "database dialect, one of [mysql,sqlite3]")
	cmd.PersistentFlags().StringVar(&root.DbConfig.DSN, "dsn", "", "database dsn(root:123456@tcp(127.0.0.1:3306)/test)")
	cmd.PersistentFlags().StringVar(&root.DbConfig.Options, "option", "", "database option(dsn?option)")
	cmd.PersistentFlags().StringVarP(&root.OutputDir, "out", "o", "./model", "out directory")
	cmd.PersistentFlags().StringSliceVarP(&root.TableNames, "table", "t", nil, "only out custom table")

	InitFlagSetForConfig(cmd.PersistentFlags(), &root.View)

	cmd.PersistentFlags().BoolVar(&root.Merge, "merge", false, "merge in a file or not")
	cmd.PersistentFlags().StringVar(&root.MergeFilename, "model", "", "merge filename")
	cmd.PersistentFlags().StringVar(&root.Template, "template", "", "use template")

	cmd.MarkPersistentFlagRequired("dsn") // nolint

	cmdAssist.Flags().StringVarP(&root.ModelPackage, "model_package", "M", "", "model package")

	cmd.AddCommand(
		cmdAssist,
		cmdMapper,
	)
	root.cmd = cmd
	return root
}
