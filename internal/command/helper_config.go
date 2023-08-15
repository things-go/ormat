package command

import (
	"github.com/spf13/pflag"
	"github.com/things-go/ens"
)

type Config struct {
	ens.Option
	DisableCommentTag bool              `yaml:"disableCommentTag" json:"disableCommentTag"` // 禁用注释放入tag标签中
	HasColumn         bool              `yaml:"hasColumn" json:"hasColumn"`                 // 是否输出字段
	SkipColumns       []string          `yaml:"skipColumns" json:"skipColumns"`             // 忽略输出字段, 格式 table.column
	Package           string            `yaml:"package" json:"package"`                     // 包名
	Options           map[string]string `yaml:"options" json:"options"`                     // 选项
	DisableField      bool              `yaml:"disableFields" json:"disableFields"`         // disable assist fields
}

func InitFlagSetForConfig(s *pflag.FlagSet, cc *Config) {
	s.StringToStringVarP(&cc.Tags, "tags", "K", map[string]string{"json": ens.TagSnakeCase}, "tags标签,类型支持[smallCamelCase,camelCase,snakeCase,kebab]")
	s.BoolVarP(&cc.EnableInt, "enableInt", "e", false, "使能int8,uint8,int16,uint16,int32,uint32输出为int,uint")
	s.BoolVarP(&cc.EnableIntegerInt, "enableIntegerInt", "E", false, "使能int32,uint32输出为int,uint")
	s.BoolVarP(&cc.EnableBoolInt, "enableBoolInt", "b", false, "使能bool输出int")
	s.BoolVarP(&cc.DisableNullToPoint, "disableNullToPoint", "B", false, "禁用字段为null时输出指针类型,将输出为sql.Nullxx")
	s.BoolVarP(&cc.DisableCommentTag, "disableCommentTag", "j", false, "禁用注释放入tag标签中")
	s.BoolVarP(&cc.EnableForeignKey, "enableForeignKey", "J", false, "使用外键")
	s.BoolVar(&cc.HasColumn, "hasColumn", false, "是否输出字段")
	s.StringSliceVar(&cc.SkipColumns, "skipColumns", nil, "忽略输出字段(仅 hasColumn = true 有效), 格式 table.column(只作用于指定表字段) 或  column(作用于所有表)")
	s.StringVar(&cc.Package, "package", "", "package name")
	s.StringToStringVar(&cc.Options, "options", nil, "options key value")

	s.BoolVar(&cc.EnableGogo, "enableGogo", false, "使能用 gogo proto (仅输出 proto 有效)")
	s.BoolVar(&cc.EnableSea, "enableSea", false, "使能用 seaql (仅输出 proto 有效)")
	s.BoolVar(&cc.DisableField, "disableField", false, "无需field字段(仅assist有效)")
}
