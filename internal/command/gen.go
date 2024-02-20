package command

import (
	"context"

	"ariga.io/atlas/sql/schema"
	"github.com/spf13/cobra"
	"github.com/things-go/ens"
	"github.com/things-go/ens/driver"
)

type genOpt struct {
	URL     string
	Tables  []string
	Exclude []string

	genFileOpt
}

type genCmd struct {
	cmd *cobra.Command
	genOpt
}

func newGenCmd() *genCmd {
	root := &genCmd{}

	getSchema := func() (ens.Schemaer, error) {
		d, err := LoadDriver(root.URL)
		if err != nil {
			return nil, err
		}
		return d.InspectSchema(context.Background(), &driver.InspectOption{
			URL: root.URL,
			InspectOptions: schema.InspectOptions{
				Mode:    schema.InspectTables,
				Tables:  root.Tables,
				Exclude: root.Exclude,
			},
		})
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

	cmdRapier := &cobra.Command{
		Use:     "rapier",
		Short:   "model rapier from database",
		Example: "ormat gen rapier",
		RunE: func(*cobra.Command, []string) error {
			sc, err := getSchema()
			if err != nil {
				return err
			}
			return root.genFileOpt.GenRapier(sc)
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

	cmd.PersistentFlags().StringVar(&root.URL, "url", "", "mysql://root:123456@127.0.0.1:3306/test")
	cmd.PersistentFlags().StringSliceVarP(&root.Tables, "table", "t", nil, "only out custom table")
	cmd.PersistentFlags().StringSliceVar(&root.Exclude, "exclude", nil, "exclude table pattern")
	cmd.PersistentFlags().StringVarP(&root.OutputDir, "out", "o", "./model", "out directory")

	InitFlagSetForConfig(cmd.PersistentFlags(), &root.View)

	cmd.PersistentFlags().BoolVar(&root.Merge, "merge", false, "merge in a file or not")
	cmd.PersistentFlags().StringVar(&root.MergeFilename, "model", "", "merge filename")
	cmd.PersistentFlags().StringVar(&root.Template, "template", "", "use template")

	cmd.MarkPersistentFlagRequired("url") // nolint

	cmdRapier.Flags().StringVarP(&root.ModelImportPath, "model_import_path", "M", "", "model import path")

	cmd.AddCommand(
		cmdRapier,
		cmdMapper,
	)
	root.cmd = cmd
	return root
}
