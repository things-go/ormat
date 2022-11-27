# ormat

database to golang struct

[![GoDoc](https://godoc.org/github.com/things-go/ormat?status.svg)](https://godoc.org/github.com/things-go/ormat)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/things-go/ormat?tab=doc)
[![codecov](https://codecov.io/gh/things-go/ormat/branch/main/graph/badge.svg)](https://codecov.io/gh/things-go/ormat)
![Action Status](https://github.com/things-go/ormat/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/things-go/ormat)](https://goreportcard.com/report/github.com/things-go/ormat)
[![Licence](https://img.shields.io/github/license/things-go/ormat)](https://raw.githubusercontent.com/things-go/ormat/main/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/things-go/ormat)](https://github.com/things-go/ormat/tags)


## Features


## Usage

### Installation

Use go get.
```bash
    go get github.com/things-go/ormat/cmd/ormat
```
or 
```bash
    go install github.com/things-go/ormat/cmd/ormat@latest
```
Example.

NOTE: 

- database filed comment `[@jsontag: realjsontag]` will overwrite the filed json tags.
- database filed comment `[@affix]` will append `,string` to the filed json tags.
- database filed comment `[@enum: enumValue]` or `[@status: enumValue]` will generate protobuf enum comment. `enumValue` should be like `{"0":["name","remark","comment"]}`
```go
// SysUser 用户表
type SysUser struct {
	ID        int64     `gorm:"column:id;type:bigint;autoIncrement;not null;primaryKey,priority:1" json:"id,omitempty"`
	Username  string    `gorm:"column:username;type:varchar(64);not null;primaryKey,priority:2;uniqueIndex:uk_username" json:"username,omitempty"`
	Password  string    `gorm:"column:password;type:varchar(255);not null" json:"password,omitempty"`
	Nickname  string    `gorm:"column:nickname;type:varchar(64);not null" json:"nickname,omitempty"`
	Phone     string    `gorm:"column:phone;type:varchar(16);not null" json:"phone,omitempty"`
	Avatar    string    `gorm:"column:avatar;type:varchar(255);not null" json:"avatar,omitempty"`
	Sex       int8      `gorm:"column:sex;type:tinyint;not null;default:3" json:"sex,omitempty"`
	Email     string    `gorm:"column:email;type:varchar(32);not null" json:"email,omitempty"`
	Status    string    `gorm:"column:status;type:varchar(1);not null;default:1" json:"status,omitempty"`
	Remark    string    `gorm:"column:remark;type:varchar(255);not null" json:"remark,omitempty"`
	Creator   string    `gorm:"column:creator;type:varchar(32);not null" json:"creator,omitempty"`
	Updator   string    `gorm:"column:updator;type:varchar(32);not null" json:"updator,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime(3);not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime(3);not null" json:"updated_at,omitempty"`
}

// TableName implement schema.Tabler interface
func (*SysUser) TableName() string {
	return "sys_user"
}
```

### Help

```shell
$ ./ormat --help

database/sql to golang struct

Usage:
  ormat [command]

Available Commands:
  build       Generate model from sql
  completion  Generate the autocompletion script for the specified shell
  config      Show/Generate config file
  expand      expand annotation from comment
  gen         Generate model/proto from database
  help        Help about any command
  sql         Generate create table sql

Flags:
  -c, --config string   config file
  -h, --help            help for ormat
  -v, --version         version for ormat

Use "ormat [command] --help" for more information about a command.
```

### Build

```shell
make linux
make windows
make mac
```
or 
```
goreleaser release --snapshot  --rm-dist
```

## References

### JetBrains OS licenses
ormat had been being developed with GoLand under the free JetBrains Open Source license(s) granted by JetBrains s.r.o., hence I would like to express my thanks here.

<a href="https://www.jetbrains.com/?from=things-go/go-modbus" target="_blank"><img src="https://github.com/thinkgos/thinkgos/blob/master/asserts/jetbrains-variant-4.svg" width="200" align="middle"/></a>


## License

This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.