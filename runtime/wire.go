//go:generate wire ./...
//go:build wireinject
// +build wireinject

package runtime

import (
	"github.com/google/wire"

	"github.com/things-go/ormat/pkg/config"
)

func NewRuntime(c *config.Config) (*Runtime, error) {
	wire.Build(
		DbSet,
		wire.NewSet(wire.Struct(new(Runtime), "*")),
	)
	return new(Runtime), nil
}
