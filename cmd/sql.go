package cmd

import (
	"github.com/spf13/cobra"

	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/pkg/database"
	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/runtime"
	"github.com/things-go/ormat/view"
)

var sqlCmd = &cobra.Command{
	Use:     "sql",
	Short:   "Generate create table sql",
	Example: "ormat sql",
	RunE: func(*cobra.Command, []string) error {
		c := config.Global
		err := c.Load()
		if err != nil {
			return err
		}
		rt, err := runtime.NewRuntime(c)
		if err != nil {
			return err
		}
		defer database.Close(rt.DB)

		vw := view.New(GetViewModel(rt), c.View)

		content, err := vw.GetDBCreateTableSQLContent()
		if err != nil {
			return err
		}
		_ = utils.WriteFile(c.OutDir+"/create_table.sql", content)

		return nil
	},
}
