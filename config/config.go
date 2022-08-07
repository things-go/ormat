package config

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/things-go/ormat/view"
)

// Database connect information
type Database struct {
	Dialect  string `yaml:"dialect" json:"dialect" binding:"required,oneof=mysql sqlite3"` // mysql, sqlite3
	Host     string `yaml:"host" json:"host"`                                              // Host. 地址
	Port     int    `yaml:"port" json:"port"`                                              // Port 端口号
	Username string `yaml:"username" json:"username"`                                      // Username 用户名
	Password string `yaml:"password" json:"password"`                                      // Password 密码
	Db       string `yaml:"db" json:"db" binding:"required"`                               // Database 数据库名
	Options  string `yaml:"options" json:"options"`                                        // Options ?号后面, 如果为空, 则为 charset=utf8&parseTime=True&loc=Local&interpolateParams=True
}

// Config custom config
type Config struct {
	Deploy     string            `yaml:"deploy" json:"deploy" binding:"oneof=local dev debug uat prod"` // 布署环境
	Database   Database          `yaml:"database" json:"database"`                                      // 数据库连接信息
	OutDir     string            `yaml:"outDir" json:"outDir" binding:"required"`                       // 文件输出路径
	TypeDefine map[string]string `yaml:"typeDefine" json:"typeDefine"`                                  // 自定义数据类型
	TableNames []string          `yaml:"tableNames" json:"tableNames"`                                  // 指定输出表
	View       view.Config       `yaml:"view" json:"view"`
}

// GetDatabase Get database configuration.
func (c *Config) GetDatabase() *Database { return &c.Database }

// GetTypeDefine 获取自定义字段映射.
func (c *Config) GetTypeDefine() map[string]string { return c.TypeDefine }

func (c *Config) GetDbDSNAndDbName() (dsn, db string, err error) {
	cc := c.Database
	switch cc.Dialect {
	case "mysql":
		if cc.Options == "" {
			cc.Options = "charset=utf8&parseTime=True&loc=Local&interpolateParams=True"
		}
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
			cc.Username, cc.Password, cc.Host, cc.Port, cc.Db, cc.Options), cc.Db, nil
	case "sqlite3":
		_, dbName := filepath.Split(cc.Db)
		if dbName != "" {
			return cc.Db, dbName, nil
		}
		err = errors.New("empty sqlite3 db name")
	default:
		err = errors.New("database not found, please check database.dialect (mysql, sqlite3, mssql)")
	}
	return "", "", err
}
