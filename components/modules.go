package components

import (
	"github.com/wiktortr/go-nimble/components/autostart"
	"github.com/wiktortr/go-nimble/components/seda"
	"github.com/wiktortr/go-nimble/components/timer"
	"go.uber.org/fx"
)

var Core = fx.Module(
	"nimble-core",
	autostart.Module,
	seda.Module,
	timer.Module,
)
