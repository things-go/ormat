package ast

// File a file
type File struct {
	Version     string
	Filename    string              // file name
	PackageName string              // package name
	Imports     map[string]struct{} // import package
	Structs     []*Struct           // struct list in file
}

func (p *File) GetProtobufEnums() []*ProtobufEnum {
	enums := make([]*ProtobufEnum, 0, 64)
	for _, s := range p.Structs {
		enums = append(enums, s.ProtoEnum...)
	}
	return enums
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
