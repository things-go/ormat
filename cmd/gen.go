package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xwb1989/sqlparser"
)

var genCmd = &cobra.Command{
	Use:     "gen",
	Short:   "gen model from sql",
	Example: "ormat gen",
	RunE: func(*cobra.Command, []string) error {
		sql :=
			"CREATE TABLE `announce` (" +
				"`id` bigint NOT NULL AUTO_INCREMENT," +
				"`title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '标题'," +
				"`content` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '内容'," +
				"`priority` tinyint unsigned NOT NULL DEFAULT '255' COMMENT '优先级'," +
				"`visible` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否显示'," +
				"`created_at` datetime NOT NULL COMMENT '发布时间'," +
				"`updated_at` datetime NOT NULL," +
				"`deleted_at` bigint NOT NULL DEFAULT '0'," +
				"PRIMARY KEY (`id`)" +
				")ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='公告-面向所有人的消息';"

		stmt, err := sqlparser.Parse(sql)
		if err != nil {
			return err
			// Do something with the error
		}

		// Otherwise do something with stmt
		switch stmt := stmt.(type) {
		case *sqlparser.DDL:
			fmt.Printf("action: %#v\n", stmt.Action)
			fmt.Printf("new name: %#v\n", stmt.NewName)
			fmt.Printf("table: %#v\n", stmt.Table)
			fmt.Printf("\n Colums: \n")
			for _, column := range stmt.TableSpec.Columns {
				fmt.Printf("%#v\n", column)
			}

			fmt.Printf("\n Indexes: \n")
			for _, idx := range stmt.TableSpec.Indexes {
				fmt.Printf("%#v\n", idx)
			}

		default:
			return errors.New("sql is not ddl")
		}
		return nil

	},
}
