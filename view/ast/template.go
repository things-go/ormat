package ast

import (
	"github.com/things-go/ormat/pkg/tpl"
)

var TableNameTpl = tpl.Template.Lookup(tpl.TableName)
var ColumnNameTpl = tpl.Template.Lookup(tpl.ColumnName)
var ProtobufCommentTpl = tpl.Template.Lookup(tpl.ProtobufComment)
var ProtobufEnumTpl = tpl.Template.Lookup(tpl.ProtobufEnum)
var ProtobufEnumMappingTpl = tpl.Template.Lookup(tpl.ProtobufEnumMapping)
