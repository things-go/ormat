package cmd

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/alecthomas/chroma/quick"
	"github.com/spf13/cobra"

	"github.com/things-go/ormat/pkg/matcher"
	"github.com/things-go/ormat/view/ast"
)

type expandCmd struct {
	cmd          *cobra.Command
	inputComment string
}

func newExpandCmd() *expandCmd {
	root := &expandCmd{}
	cmd := &cobra.Command{
		Use:     "expand",
		Short:   "Expand annotation from comment",
		Example: "ormat expand -i comment",
		RunE: func(*cobra.Command, []string) error {
			str := matcher.EnumAnnotation(root.inputComment)
			if str == "" {
				return errors.New("没有符合的注解")
			}
			mp, err := ast.ParseEnumAnnotation(str)
			if err != nil {
				return err
			}
			v, err := json.MarshalIndent(mp, " ", "  ")
			if err != nil {
				return err
			}
			return quick.Highlight(os.Stdout, string(v), "JSON", "terminal", "solarized-dark")
		},
	}
	cmd.Flags().StringVarP(&root.inputComment, "input", "i", "", "input comment")
	cmd.MarkFlagRequired("input") // nolint

	root.cmd = cmd
	return root
}
