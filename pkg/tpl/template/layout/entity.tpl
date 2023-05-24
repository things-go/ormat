// Code generated by ormat. DO NOT EDIT.
// version: {{.Version}}

package {{.PackageName}}

import (
    "gorm.io/gorm"
)

{{- range $e := .Structs}}

type {{$e.StructName}}_Entity struct {
	db *gorm.DB
}

type {{$e.StructName}}_Executor struct {
	db *gorm.DB
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
		funcs: make([]func(*gorm.DB) *gorm.DB, 0, 16),
	}
}

// PreProcess on db
func (x *{{$e.StructName}}_Executor) PreProcess(funcs ...func(*gorm.DB) *gorm.DB) *{{$e.StructName}}_Executor {
	db := x.db
	for _, f := range funcs {
		db = f(db)
	}
	x.db = db
	return x
}

// Condition additional condition
func (x *{{$e.StructName}}_Executor) Condition(funcs ...func(*gorm.DB) *gorm.DB) *{{$e.StructName}}_Executor {
	x.funcs = append(x.funcs, funcs...)
	return x
}

// Condition additional condition to executor
func (x *{{$e.StructName}}_Executor) Where(query any, args ...any) *{{$e.StructName}}_Executor {
	f := func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	}
	x.funcs = append(x.funcs, f)
	return x
}

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
	err = x.db.Model(&{{$e.StructName}}{}).
		Scopes(x.funcs...).
		Count(&count).Error
	return count, err
}

func (x *{{$e.StructName}}_Executor) First(dest any) error {
	return x.db.Model(&{{$e.StructName}}{}).
		Scopes(x.funcs...).
		First(dest).Error
}

func (x *{{$e.StructName}}_Executor) Take(dest any) error {
	return x.db.Model(&{{$e.StructName}}{}).
		Scopes(x.funcs...).
		Take(dest).Error
}

func (x *{{$e.StructName}}_Executor) Last(dest any) error {
	return x.db.Model(&{{$e.StructName}}{}).
		Scopes(x.funcs...).
		Last(dest).Error
}

func (x *{{$e.StructName}}_Executor) Scan(dest any) error {
	return x.db.Model(&{{$e.StructName}}{}).
		Scopes(x.funcs...).
		Scan(dest).Error
}

func (x *{{$e.StructName}}_Executor) Pluck(column string, value any) error {
	return x.db.Model(&{{$e.StructName}}{}).
		Scopes(x.funcs...).
		Pluck(column, value).Error
}

func (x *{{$e.StructName}}_Executor) Exist() (exist bool, err error) {
	err = x.db.Model(&{{$e.StructName}}{}).
		Select("1").
		Scopes(x.funcs...).
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
	return x.db.Model(&{{$e.StructName}}{}).
		Scopes(x.funcs...).
		Find(dest).Error
}

func (x *{{$e.StructName}}_Executor) Create(value any) error {
	return x.db.Scopes(x.funcs...).Create(value).Error
}

func (x *{{$e.StructName}}_Executor) CreateInBatches(value any, batchSize int) error {
	return x.db.Scopes(x.funcs...).CreateInBatches(value, batchSize).Error
}

func (x *{{$e.StructName}}_Executor) Save(value any) error {
	return x.db.Model(&{{$e.StructName}}{}).
		Scopes(x.funcs...).
		Save(value).Error
}

func (x *{{$e.StructName}}_Executor) Updates(value *{{$e.StructName}}) error {
	return x.db.Model(&{{$e.StructName}}{}).
		Scopes(x.funcs...).
		Updates(value).Error
}

func (x *{{$e.StructName}}_Executor) Update(column string, value any) error {
	return x.db.Model(&{{$e.StructName}}{}).
		Scopes(x.funcs...).
		Update(column, value).Error
}

func (x *{{$e.StructName}}_Executor) UpdateColumns(value *{{$e.StructName}}) error {
	return x.db.Model(&{{$e.StructName}}{}).
		Scopes(x.funcs...).
		UpdateColumns(value).Error
}

func (x *{{$e.StructName}}_Executor) UpdateColumn(column string, value any) error {
	return x.db.Model(&{{$e.StructName}}{}).
		Scopes(x.funcs...).
		UpdateColumn(column, value).Error
}

func (x *{{$e.StructName}}_Executor) Delete() error {
	return x.db.Model(&{{$e.StructName}}{}).
		Scopes(x.funcs...).
		Delete(&{{$e.StructName}}{}).Error
}

{{- end}}