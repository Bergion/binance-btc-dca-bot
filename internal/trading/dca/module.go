package dca

import (
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewExecutor,
)
