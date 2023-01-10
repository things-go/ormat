package ast

// File a file
type File struct {
	Version     string
	Filename    string              // file name
	PackageName string              // package name
	Imports     map[string]struct{} // import package
	Structs     []*Struct           // struct list in file
	HasColumn   bool                // has column
	HasHelper   bool                // has helper
}

func IntoImports(s []*Struct) map[string]struct{} {
	mp := make(map[string]struct{})
	for _, v := range s {
		for _, v1 := range v.StructFields {
			if v2, ok := ImportsHeads[v1.FieldType]; ok && v2 != "" {
				mp[v2] = struct{}{}
			}
		}
	}
	return mp
}
