package ast

type ProtobufEnumFile struct {
	Version string
	Package string
	Options map[string]string
	Enums   []*ProtobufEnum
}
