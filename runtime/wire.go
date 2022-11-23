//go:generate wire ./...
//go:build wireinject
// +build wireinject

package runtime

import (
	"github.com/google/wire"
)

func NewRuntime(remote bool) (*Runtime, error) {
	wire.Build(
		ConfigSet,
		DbSet,
		wire.NewSet(wire.Struct(new(Runtime), "*")),
	)
	return new(Runtime), nil
}
