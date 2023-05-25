// Code generated by ormat. DO NOT EDIT.
// version: {{.Version}}

package {{.PackageName}}

import (
	"context"

	assist "github.com/things-go/gorm-assist"
    "gorm.io/gorm"
	"gorm.io/gorm/clause"
)

{{- range $e := .Structs}}

type {{$e.StructName}}_Entity struct {
	db *gorm.DB
}

type {{$e.StructName}}_Executor struct {
	db *gorm.DB
	table func(*gorm.DB) *gorm.DB
	funcs []func(*gorm.DB) *gorm.DB
}

func New_{{$e.StructName}}(db *gorm.DB) *{{$e.StructName}}_Entity {
	return &{{$e.StructName}}_Entity{
		db: db,
	}
}

// Executor new executor
func (x *{{$e.StructName}}_Entity) Executor() *{{$e.StructName}}_Executor {
	return &{{$e.StructName}}_Executor{
		db: x.db,
		table: nil,
		funcs: make([]func(*gorm.DB) *gorm.DB, 0, 16),
	}
}

func (x *{{$e.StructName}}_Executor) Session(config *gorm.Session) *{{$e.StructName}}_Executor {
	x.db = x.db.Session(config)
	return x
}

func (x *{{$e.StructName}}_Executor) WithContext(ctx context.Context) *{{$e.StructName}}_Executor {
	x.db = x.db.WithContext(ctx)
	return x
}

func (x *{{$e.StructName}}_Executor) Debug() *{{$e.StructName}}_Executor {
	x.db = x.db.Debug()
	return x
}

/********************************** chains api *********************************/

func (x *{{$e.StructName}}_Executor) Clauses(conds ...clause.Expression) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, func(db *gorm.DB) *gorm.DB {
		return db.Clauses(conds...)
	})
	return x
}

func (x *{{$e.StructName}}_Executor) Table(name string, args ...any) *{{$e.StructName}}_Executor {
	x.table = func(db *gorm.DB) *gorm.DB {
		return db.Table(name, args...)
	}
	return x
}

func (x *{{$e.StructName}}_Executor) Distinct(args ...any) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, func(db *gorm.DB) *gorm.DB {
		return db.Distinct(args...)
	})
	return x
}

func (x *{{$e.StructName}}_Executor) Select(query any, args ...any) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, func(db *gorm.DB) *gorm.DB {
		return db.Select(query, args...)
	})
	return x
}

func (x *{{$e.StructName}}_Executor) Omit(columns ...string) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, func(db *gorm.DB) *gorm.DB {
		return db.Omit(columns...)
	})
	return x
}

func (x *{{$e.StructName}}_Executor) Where(query any, args ...any) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	})
	return x
}

func (x *{{$e.StructName}}_Executor) Not(query any, args ...any) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, func(db *gorm.DB) *gorm.DB {
		return db.Not(query, args...)
	})
	return x
}

func (x *{{$e.StructName}}_Executor) Or(query any, args ...any) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, func(db *gorm.DB) *gorm.DB {
		return db.Or(query, args...)
	})
	return x
}

func (x *{{$e.StructName}}_Executor) Joins(query string, args ...any) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, func(db *gorm.DB) *gorm.DB {
		return db.Joins(query, args...)
	})
	return x
}

func (x *{{$e.StructName}}_Executor) InnerJoins(query string, args ...any) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, func(db *gorm.DB) *gorm.DB {
		return db.InnerJoins(query, args...)
	})
	return x
}

func (x *{{$e.StructName}}_Executor) Group(name string) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, func(db *gorm.DB) *gorm.DB {
		return db.Group(name)
	})
	return x
}

func (x *{{$e.StructName}}_Executor) Having(query any, args ...any) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, func(db *gorm.DB) *gorm.DB {
		return db.Having(query, args...)
	})
	return x
}

func (x *{{$e.StructName}}_Executor) Order(value any) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, func(db *gorm.DB) *gorm.DB {
		return db.Order(value)
	})
	return x
}

func (x *{{$e.StructName}}_Executor) Limit(limit int) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	})
	return x
}

func (x *{{$e.StructName}}_Executor) Offset(offset int) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	})
	return x
}

func (x *{{$e.StructName}}_Executor) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, funcs...)
	return x
}

func (x *{{$e.StructName}}_Executor) TableExpr(fromSubs ...assist.From) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, assist.Table(fromSubs...))
	return x
}

func (x *{{$e.StructName}}_Executor) SelectExpr(columns ...assist.Expr) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, assist.Select(columns...))
	return x
}

func (x *{{$e.StructName}}_Executor) OrderExpr(columns ...assist.Expr) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, assist.Order(columns...))
	return x
}

func (x *{{$e.StructName}}_Executor) GroupExpr(columns ...assist.Expr) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, assist.Group(columns...))
	return x
}

func (x *{{$e.StructName}}_Executor) CrossJoinsExpr(tableName string, conds ...assist.Expr) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, assist.CrossJoins(tableName, conds...))
	return x
}

func (x *{{$e.StructName}}_Executor) CrossJoinsXExpr(tableName, alias string, conds ...assist.Expr) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, assist.CrossJoinsX(tableName, alias, conds...))
	return x
}


func (x *{{$e.StructName}}_Executor) InnerJoinsExpr(tableName string, conds ...assist.Expr) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, assist.InnerJoins(tableName, conds...))
	return x
}

func (x *{{$e.StructName}}_Executor) InnerJoinsXExpr(tableName, alias string, conds ...assist.Expr) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, assist.InnerJoinsX(tableName, alias, conds...))
	return x
}

func (x *{{$e.StructName}}_Executor) LeftJoinsExpr(tableName string, conds ...assist.Expr) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, assist.LeftJoins(tableName, conds...))
	return x
}

func (x *{{$e.StructName}}_Executor) LeftJoinsXExpr(tableName, alias string, conds ...assist.Expr) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, assist.LeftJoinsX(tableName, alias, conds...))
	return x
}

func (x *{{$e.StructName}}_Executor) RightJoinsExpr(tableName string, conds ...assist.Expr) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, assist.RightJoins(tableName, conds...))
	return x
}

func (x *{{$e.StructName}}_Executor) RightJoinsXExpr(tableName, alias string, conds ...assist.Expr) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, assist.RightJoinsX(tableName, alias, conds...))
	return x
}

func (x *{{$e.StructName}}_Executor) LockingUpdate() *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, assist.LockingUpdate())
	return x
}

func (x *{{$e.StructName}}_Executor) LockingShare() *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, assist.LockingShare())
	return x
}

func (x *{{$e.StructName}}_Executor) Pagination(page, perPage int64, maxPerPages ...int64) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, assist.Pagination(page, perPage, maxPerPages...))
	return x
}

func (x *{{$e.StructName}}_Executor) chains() (db *gorm.DB) {
	if x.table == nil {
		db = x.db.Model(&{{$e.StructName}}{})
	} else {
		db = x.db.Scopes(x.table)
	}
	return db.Scopes(x.funcs...)
}

/********************************** finish api *********************************/

func (x *{{$e.StructName}}_Executor) FirstOne() (*{{$e.StructName}}, error) {
	var row {{$e.StructName}}

	err := x.First(&row)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (x *{{$e.StructName}}_Executor) TakeOne() (*{{$e.StructName}}, error) {
	var row {{$e.StructName}}

	err := x.Take(&row)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (x *{{$e.StructName}}_Executor) LastOne() (*{{$e.StructName}}, error) {
	var row {{$e.StructName}}

	err := x.Last(&row)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (x *{{$e.StructName}}_Executor) ScanOne() (*{{$e.StructName}}, error) {
	var row {{$e.StructName}}

	err := x.Scan(&row)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (x *{{$e.StructName}}_Executor) Count() (count int64, err error) {
	err = x.chains().Count(&count).Error
	return count, err
}

func (x *{{$e.StructName}}_Executor) First(dest any) error {
	return x.chains().First(dest).Error
}

func (x *{{$e.StructName}}_Executor) Take(dest any) error {
	return x.chains().Take(dest).Error
}

func (x *{{$e.StructName}}_Executor) Last(dest any) error {
	return x.chains().Last(dest).Error
}

func (x *{{$e.StructName}}_Executor) Scan(dest any) error {
	return x.chains().Scan(dest).Error
}

func (x *{{$e.StructName}}_Executor) Pluck(column string, value any) error {
	return x.chains().Pluck(column, value).Error
}

func (x *{{$e.StructName}}_Executor) Exist() (exist bool, err error) {
	err = x.chains().
		Select("1").
		Limit(1).
		Scan(&exist).Error
	return exist, err
}

func (x *{{$e.StructName}}_Executor) FindAll() ([]*{{$e.StructName}}, error) {
	var rows []*{{$e.StructName}}

	err := x.Find(&rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (x *{{$e.StructName}}_Executor) Find(dest any) error {
	return x.chains().Find(dest).Error
}

func (x *{{$e.StructName}}_Executor) Create(value any) error {
	return x.db.Scopes(x.funcs...).Create(value).Error
}

func (x *{{$e.StructName}}_Executor) CreateInBatches(value any, batchSize int) error {
	return x.db.Scopes(x.funcs...).CreateInBatches(value, batchSize).Error
}

func (x *{{$e.StructName}}_Executor) Save(value any) error {
	return x.db.Scopes(x.funcs...).Save(value).Error
}

func (x *{{$e.StructName}}_Executor) Updates(value *{{$e.StructName}}) error {
	return x.chains().Updates(value).Error
}

func (x *{{$e.StructName}}_Executor) Update(column string, value any) error {
	return x.chains().Update(column, value).Error
}

func (x *{{$e.StructName}}_Executor) UpdateColumns(value *{{$e.StructName}}) error {
	return x.chains().UpdateColumns(value).Error
}

func (x *{{$e.StructName}}_Executor) UpdateColumn(column string, value any) error {
	return x.chains().UpdateColumn(column, value).Error
}

func (x *{{$e.StructName}}_Executor) Delete() error {
	return x.chains().Delete(&{{$e.StructName}}{}).Error
}

{{- end}}