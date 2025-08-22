package nimble

import "go.uber.org/fx"

type ComponentImpl interface {
	Start()
	Stop()
	Inbound() chan *Message
	Process(msg *Message) error
}

type Component interface {
	Key() string
	Instantiate(reg Registry, params *ComponentParams) (ComponentImpl, error)
}

func AsComponent(f any) fx.Option {
	return fx.Provide(fx.Annotate(
		f,
		fx.ResultTags(`group:"nimble-components"`),
	))
}
