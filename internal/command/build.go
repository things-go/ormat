package command

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/things-go/ens"
	driverMysql "github.com/things-go/ens/driver/mysql"
	"golang.org/x/exp/slog"
)

type buildOpt struct {
	InputFile []string

	genFileOpt
}

type buildCmd struct {
	cmd *cobra.Command
	buildOpt
}

func newBuildCmd() *buildCmd {
	root := &buildCmd{}

	getSchema := func() ens.Schemaer {
		innerParseFromFile := func(filename string) (ens.Schemaer, error) {
			content, err := os.ReadFile(filename)
			if err != nil {
				return nil, err
			}
			d := &driverMysql.SQL{
				CreateTableSQL:    string(content),
				DisableCommentTag: false,
			}
			return d.GetSchema()
		}

		mixin := &ens.MixinSchema{
			Name:     "",
			Entities: make([]ens.MixinEntity, 0, 128),
		}
		for _, filename := range root.InputFile {
			sc, err := innerParseFromFile(filename)
			if err != nil {
				slog.Warn("🧐 parse from SQL file(%s) failed !!!", filename)
				continue
			}
			mixin.Entities = append(mixin.Entities, sc.(*ens.MixinSchema).Entities...)
		}
		return mixin
	}

	cmd := &cobra.Command{
		Use:     "build",
		Short:   "Generate model from sql",
		Example: "ormat build",
		RunE: func(*cobra.Command, []string) error {
			sc := getSchema()
			return root.genFileOpt.GenModel(sc)
		},
	}

	cmdAssist := &cobra.Command{
		Use:     "assist",
		Short:   "model assist from sql",
		Example: "ormat build assist",
		RunE: func(*cobra.Command, []string) error {
			sc := getSchema()
			return root.genFileOpt.GenAssist(sc)
		},
	}

	cmdMapper := &cobra.Command{
		Use:     "mapper",
		Short:   "model mapper from database",
		Example: "ormat gen mapper",
		RunE: func(*cobra.Command, []string) error {
			sc := getSchema()
			return root.genFileOpt.GenMapper(sc)
		},
	}

	cmd.PersistentFlags().StringSliceVarP(&root.InputFile, "input", "i", nil, "input file")
	cmd.PersistentFlags().StringVarP(&root.OutputDir, "out", "o", "./model", "out directory")

	InitFlagSetForConfig(cmd.PersistentFlags(), &root.View)

	cmd.PersistentFlags().BoolVar(&root.Merge, "merge", false, "merge in a file or not")
	cmd.PersistentFlags().StringVar(&root.MergeFilename, "filename", "", "merge filename")
	cmd.PersistentFlags().StringVar(&root.Template, "template", "__in_go", "use custom template")

	cmdAssist.Flags().StringVarP(&root.ModelPackage, "model_package", "M", "", "model package")

	cmd.MarkPersistentFlagRequired("input") // nolint
	cmd.AddCommand(
		cmdAssist,
		cmdMapper,
	)
	root.cmd = cmd
	return root
}
