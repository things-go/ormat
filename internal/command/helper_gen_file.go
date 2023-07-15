package command

import (
	"github.com/things-go/ens"
	"github.com/things-go/ens/codegen"
	"github.com/things-go/ens/utils"
	"golang.org/x/exp/slog"
	"golang.org/x/tools/imports"
)

type genFileOpt struct {
	OutputDir     string
	View          Config
	Merge         bool
	MergeFilename string
	Template      string
	// assist命令 模型包名
	ModelPackage string
}

func (self *genFileOpt) build(mixin ens.Schemaer) *ens.Schema {
	return mixin.Build(&self.View.Option)
}

func (self *genFileOpt) GenModel(mixin ens.Schemaer) error {
	skipColumns := make(map[string]struct{})
	for _, v := range self.View.SkipColumns {
		skipColumns[v] = struct{}{}
	}

	codegenOption := []codegen.Option{
		codegen.WithByName("ormat"),
		codegen.WithVersion(version),
		codegen.WithPackageName(utils.GetPkgName(self.OutputDir)),
		codegen.WithOptions(self.View.Options),
		codegen.WithSkipColumns(skipColumns),
		codegen.WithHasColumn(self.View.HasColumn),
	}
	sc := self.build(mixin)
	if self.Merge {
		data, err := codegen.New(sc.Entities, codegenOption...).GenModel().FormatSource()
		if err != nil {
			return err
		}
		filename := joinFilename(self.OutputDir, self.MergeFilename, ".go")
		err = WriteFile(filename, data)
		if err != nil {
			return err
		}
		slog.Info("👉 " + filename)
	} else {
		for _, entity := range sc.Entities {
			data, err := codegen.New([]*ens.Entity{entity}, codegenOption...).GenModel().FormatSource()
			if err != nil {
				slog.Error(err.Error(), entity.Name)
				return err
			}
			filename := joinFilename(self.OutputDir, entity.Name, ".go")
			data, err = imports.Process(filename, data, nil)
			if err != nil {
				return err
			}
			err = WriteFile(filename, data)
			if err != nil {
				return err
			}
			slog.Info("👉 " + filename)
		}
	}
	slog.Info("😄 generate success !!!")
	return nil
}

func (self *genFileOpt) GenAssist(mixin ens.Schemaer) error {
	skipColumns := make(map[string]struct{})
	for _, v := range self.View.SkipColumns {
		skipColumns[v] = struct{}{}
	}

	codegenOption := []codegen.Option{
		codegen.WithByName("ormat"),
		codegen.WithVersion(version),
		codegen.WithPackageName(utils.GetPkgName(self.OutputDir)),
		codegen.WithOptions(self.View.Options),
		codegen.WithSkipColumns(skipColumns),
		codegen.WithHasColumn(self.View.HasColumn),
	}
	sc := self.build(mixin)
	for _, entity := range sc.Entities {
		data, err := codegen.New([]*ens.Entity{entity}, codegenOption...).
			GenAssist(self.ModelPackage).
			FormatSource()
		if err != nil {
			return err
		}
		filename := joinFilename(self.OutputDir, entity.Name, ".assist.go")
		data, err = imports.Process(filename, data, nil)
		if err != nil {
			return err
		}
		err = WriteFile(filename, data)
		if err != nil {
			return err
		}
		slog.Info("👉 " + filename)
	}
	slog.Info("😄 generate success !!!")
	return nil
}

func (self *genFileOpt) GenMapper(mixin ens.Schemaer) error {
	skipColumns := make(map[string]struct{})
	for _, v := range self.View.SkipColumns {
		skipColumns[v] = struct{}{}
	}
	codegenOption := []codegen.Option{
		codegen.WithByName("ormat"),
		codegen.WithVersion(version),
		codegen.WithPackageName(utils.GetPkgName(self.OutputDir)),
		codegen.WithOptions(self.View.Options),
		codegen.WithSkipColumns(skipColumns),
		codegen.WithHasColumn(self.View.HasColumn),
	}
	sc := self.build(mixin)
	for _, entity := range sc.Entities {
		data := codegen.New([]*ens.Entity{entity}, codegenOption...).
			GenMapper().
			Bytes()
		filename := joinFilename(self.OutputDir, entity.Name, ".proto")
		err := WriteFile(filename, data)
		if err != nil {
			return err
		}
		slog.Info("👉 " + filename)
	}
	slog.Info("😄 generate success !!!")
	return nil
}
