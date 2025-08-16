package timer

import (
	"github.com/wiktortr/go-nimble/nimble"
	"go.uber.org/fx"
)

type Timer struct {
}

func (p *Timer) Key() string {
	return "timer"
}

func (p *Timer) Instantiate(params *nimble.ComponentParams) (nimble.ComponentImpl, error) {
	duration, err := params.Duration("interval")
	if err != nil {
		return nil, err
	}

	return &timerImpl{
		dur:     duration,
		channel: make(chan *nimble.Message),
	}, nil
}

var Module = fx.Module(
	"nimble-timer",
	nimble.AsComponent(func() nimble.Component {
		return &Timer{}
	}),
)
