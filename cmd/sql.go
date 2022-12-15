package cmd

import (
	"github.com/spf13/cobra"

	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/pkg/database"
	"github.com/things-go/ormat/pkg/tpl"
	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/runtime"
	"github.com/things-go/ormat/view"
)

type sqlCmd struct {
	cmd *cobra.Command
}

func newSqlCmd() *sqlCmd {
	root := &sqlCmd{}
	cmd := &cobra.Command{
		Use:     "sql",
		Short:   "Generate create table sql",
		Example: "ormat sql",
		RunE: func(*cobra.Command, []string) error {
			c := config.Global
			err := c.Load()
			if err != nil {
				return err
			}
			setupBase(c)
			rt, err := runtime.NewRuntime(c)
			if err != nil {
				return err
			}
			defer database.Close(rt.DB)

			vw := view.New(GetViewModel(rt), c.View)

			sqlFile, err := vw.GetDBCreateTableSQL()
			if err != nil {
				return err
			}
			return utils.WriteFileWithTemplate(c.OutDir+"/create_table.sql", tpl.SqlDDL, sqlFile)
		},
	}
	root.cmd = cmd
	return root
}
