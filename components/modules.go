package components

import (
	"github.com/wiktortr/go-nimble/components/seda"
	"github.com/wiktortr/go-nimble/components/timer"
	"go.uber.org/fx"
)

var Core = fx.Module(
	"nimble-core",
	seda.Module,
	timer.Module,
)
