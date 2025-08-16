package nimble

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MsgProcessor func(message *Message) error

type Route struct {
	id           string
	Name         string
	From         string
	Concurrency  int
	Ctx          context.Context
	cancel       context.CancelFunc
	dependencies []string
	components   map[string]ComponentImpl
	Registry     *Registry
	mainBlock    Block
	currentBlock Block
	compiledFunc MsgProcessor
}

func NewRoute(name string) *Route {
	mainBlock := &LinearBlock{nil, []Block{}}
	return &Route{
		id:           uuid.NewString(),
		Name:         name,
		Concurrency:  1,
		dependencies: []string{},
		components:   make(map[string]ComponentImpl),
		mainBlock:    mainBlock,
		currentBlock: mainBlock,
	}
}

func From(uri string) *Route {
	route := NewRoute(uri)
	route.From = uri
	route.dependencies = append(route.dependencies, uri)
	return route
}

func (m *Route) To(uri string) *Route {
	m.dependencies = append(m.dependencies, uri)
	err := m.currentBlock.AddBlock(&ComponentBlock{uri})
	if err != nil {
		m.Registry.logger.Error("Failed to add block", zap.Error(err))
	}
	return m
}

func (m *Route) Log(logMsg string) *Route {
	err := m.currentBlock.AddBlock(&FunctionalBlock{
		func(message *Message) error {
			m.Registry.logger.Info(logMsg, zap.Any("msg", message))
			return nil
		},
	})
	if err != nil {
		m.Registry.logger.Error("Failed to add block", zap.Error(err))
	}
	return m
}

func (m *Route) End() *Route {
	m.currentBlock = m.currentBlock.GetParent()
	return m
}

func (m *Route) Run(msg *Message) error {
	if m.compiledFunc == nil {
		var err error
		m.compiledFunc, err = m.mainBlock.Compile(m.Registry)
		if err != nil {
			return err
		}
	}
	return m.compiledFunc(msg)
}

func (m *Route) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	m.Ctx = ctx
	m.cancel = cancel

	for i := 0; i < m.Concurrency; i++ {
		go m.handle()
	}
}

func (m *Route) handle() {
	fromComp := m.components[m.From]

	for {
		select {
		case <-m.Ctx.Done():
			return
		case msg := <-fromComp.Inbound():
			err := m.Run(msg)
			if err != nil {
				m.Registry.logger.Error(err.Error(), zap.String("route", m.Name), zap.Any("msg", msg))
			}
		}
	}
}

func (m *Route) Stop() {
	m.cancel()
}

func AsRoute(f any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			f,
			fx.ResultTags(`group:"nimble-routes"`),
		),
	)
}
