package config

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/thinkgos/ormat/pkg/env"
	"github.com/thinkgos/ormat/pkg/infra"
	"github.com/thinkgos/ormat/pkg/zapl"
	"github.com/thinkgos/ormat/view"
)

// Config custom config struct
type Config struct {
	Deploy     string            `yaml:"deploy" binding:"required,oneof=local dev debug uat prod"` // 布署环境
	Database   Database          `yaml:"database"`                                                 // 数据库连接信息
	OutDir     string            `yaml:"outDir"`                                                   // 输出路径
	TypeDefine map[string]string `yaml:"typeDefine"`                                               // 自定义类型
	TableNames []string          `yaml:"tableNames"`                                               // 指定表
	View       view.Config
}

// Database information
type Database struct {
	Dialect  string `binding:"required,oneof=mysql sqlite3"` // mysql, sqlite3
	Host     string // Host. 地址
	Port     int    // Port 端口号
	Username string // Username 用户名
	Password string // Password 密码
	Db       string `binding:"required"` // Database 数据库名
}

var cfg = Config{
	Deploy: env.DeployProd,
	Database: Database{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "root",
		Password: "root",
		Db:       "test",
	},
	OutDir:     "./model",
	TypeDefine: make(map[string]string),
	View: view.Config{
		DbTag:         "gorm",
		WebTags:       []view.WebTag{{Kind: "snakeCase", Tag: "json", HasOmit: true}},
		DisableNull:   false,
		EnableInt:     false,
		IsNullToPoint: true,
		IsOutSQL:      false,
		IsForeignKey:  false,
		IsCommentTag:  false,
	},
}

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(infra.GetExecutableDir())
	viper.AddConfigPath(infra.GetWd())
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return viper.Unmarshal(&cfg, func(c *mapstructure.DecoderConfig) { c.TagName = "yaml" })
}

// GetDbInfo Get configuration information
func GetDbInfo() Database {
	return cfg.Database
}

// GetConfig get config
func GetConfig() Config {
	return cfg
}

// GetTypeDefine 获取自定义字段映射
func GetTypeDefine() map[string]string {
	return cfg.TypeDefine
}

func GetDatabaseDSNAndDbName() (dsn string, db string) {
	c := cfg.Database
	switch c.Dialect {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&interpolateParams=True",
			c.Username, c.Password, c.Host, c.Port, c.Db), c.Db
	case "sqlite3":
		_, dbName := filepath.Split(c.Db)
		if dbName == "" {
			panic("sqlite3: invalid db name")
		}
		return c.Db, dbName
	}
	zapl.Fatal("database not found: please check database.dialect (mysql, sqlite3, mssql)")
	return "", ""
}
