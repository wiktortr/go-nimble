package seda

import (
	"github.com/wiktortr/go-nimble/nimble"
	"go.uber.org/fx"
)

type Seda struct {
}

func (p *Seda) Key() string {
	return "seda"
}

func (p *Seda) Instantiate(_ nimble.Registry, params *nimble.ComponentParams) (nimble.ComponentImpl, error) {

	buffSize := params.IntDef("buffSize", 100)

	return &sedaImpl{make(chan *nimble.Message, buffSize)}, nil
}

var Module = fx.Module(
	"nimble-seda",
	nimble.AsComponent(func() nimble.Component {
		return &Seda{}
	}),
)
