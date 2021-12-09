package ast

import (
	"strings"
)

// File a file
type File struct {
	Name        string            // file name
	PackageName string            // package name
	Imports     map[string]string // import package
	Structs     []Struct          // struct list
}

// SetName set file name
func (p *File) SetName(name string) *File {
	p.Name = name
	return p
}

// GetName get file name
func (p *File) GetName() string { return p.Name }

// SetPackageName set file package name
func (p *File) SetPackageName(name string) *File {
	p.PackageName = name
	return p
}

// GetPackageName get package name
func (p *File) GetPackageName() string { return p.PackageName }

// AddImport Add import by type
func (p *File) AddImport(imp string) *File {
	if p.Imports == nil {
		p.Imports = make(map[string]string)
	}
	p.Imports[imp] = imp
	return p
}

// AddStruct Add a structure
func (p *File) AddStruct(st Struct) *File {
	p.Structs = append(p.Structs, st)
	return p
}

// Build Get the result data
func (p *File) Build() string {
	buf := strings.Builder{}

	buf.WriteString("package" + delimTab + p.PackageName + delimLF)

	// auto add import
	for _, v := range p.Structs {
		for _, v1 := range v.Fields {
			if v2, ok := ImportsHeads[v1.GetType()]; ok {
				if v2 != "" {
					p.AddImport(v2)
				}
			}
		}
	}

	// add imports
	if len(p.Imports) > 0 {
		buf.WriteString("import (" + delimLF)
		for _, v := range p.Imports {
			buf.WriteString(v + delimLF)
		}
		buf.WriteString(")" + delimLF)
	}

	// add struct
	for _, v := range p.Structs {
		for _, v1 := range v.BuildLines() {
			buf.WriteString(v1 + delimLF)
		}
		// add table name function
		buf.WriteString(v.BuildTableNameTemplate() + delimLF)
		buf.WriteString(delimLF)
		buf.WriteString(v.BuildColumnNameTemplate())
	}
	return buf.String()
}

// ImportsHeads import head options
var ImportsHeads = map[string]string{
	"string":         `"string"`,
	"time.Time":      `"time"`,
	"gorm.Model":     `"gorm.io/gorm"`,
	"fmt":            `"fmt"`,
	"datatypes.JSON": `"gorm.io/datatypes"`,
	"datatypes.Date": `"gorm.io/datatypes"`,
}
