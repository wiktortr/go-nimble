package autostart

import (
	"github.com/wiktortr/go-nimble/nimble"
	"go.uber.org/fx"
	"time"
)

type AutoStart struct {
}

func (p *AutoStart) Key() string {
	return "autostart"
}

func (p *AutoStart) Instantiate(_ nimble.Registry, params *nimble.ComponentParams) (nimble.ComponentImpl, error) {

	delay := params.DurationDef("delay", time.Second)

	return &autostartImpl{
		delay:   delay,
		channel: make(chan *nimble.Message, 1),
	}, nil
}

var Module = fx.Module(
	"nimble-autostart",
	nimble.AsComponent(func() nimble.Component {
		return &AutoStart{}
	}),
)
