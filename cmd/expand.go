package cmd

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/alecthomas/chroma/quick"
	"github.com/spf13/cobra"

	"github.com/things-go/ormat/view/ast"
)

var inputComment string

func init() {
	expandCmd.Flags().StringVarP(&inputComment, "input", "i", "", "input file")
	expandCmd.MarkFlagRequired("input") // nolint
}

var expandCmd = &cobra.Command{
	Use:     "expand",
	Short:   "expand annotation from comment",
	Example: "ormat expand -i comment",
	RunE: func(*cobra.Command, []string) error {
		str := ast.MatchEnumAnnotation(inputComment)
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
