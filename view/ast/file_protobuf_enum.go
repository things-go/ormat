package ast

type ProtobufEnumFile struct {
	Version     string
	PackageName string
	Package     string
	Options     map[string]string
	Enums       []*ProtobufEnum
}
