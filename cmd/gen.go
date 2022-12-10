package cmd

import (
	"github.com/spf13/cobra"

	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/pkg/database"
	"github.com/things-go/ormat/pkg/utils"
	"github.com/things-go/ormat/runtime"
	"github.com/things-go/ormat/view"
)

var genCmd = &cobra.Command{
	Use:     "gen",
	Short:   "Generate model/proto from database",
	Example: "ormat gen",
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

		astFiles, err := view.New(GetViewModel(rt), c.View).
			GetDbFile(utils.GetPkgName(c.OutDir))
		if err != nil {
			return err
		}
		genModelFile(astFiles, &c.View.Protobuf, c.OutDir)
		return nil
	},
}
