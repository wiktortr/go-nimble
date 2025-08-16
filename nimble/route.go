package nimble

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"log"
)

type MsgProcessor func(message *Message) error

type Route struct {
	id           string
	Name         string
	Concurrency  int
	Ctx          context.Context
	cancel       context.CancelFunc
	From         string
	dependencies []string
	steps        []MsgProcessor
	components   map[string]ComponentImpl
	Registry     *Registry
}

func From(uri string) *Route {
	return &Route{
		id:           uuid.NewString(),
		Name:         uri,
		From:         uri,
		Concurrency:  1,
		dependencies: []string{uri},
		steps:        []MsgProcessor{},
		components:   make(map[string]ComponentImpl),
	}
}

func (route *Route) Process(process MsgProcessor) *Route {
	route.steps = append(route.steps, process)
	return route
}

func (route *Route) To(uri string) *Route {
	route.dependencies = append(route.dependencies, uri)
	route.steps = append(route.steps, func(message *Message) error {
		return route.components[uri].Process(message)
	})
	return route
}

func (route *Route) Log(logMsg string) *Route {
	route.steps = append(route.steps, func(message *Message) error {
		route.Registry.logger.Info(logMsg, zap.Any("msg", message))
		return nil
	})
	return route
}

func (route *Route) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	route.Ctx = ctx
	route.cancel = cancel

	for i := 0; i < route.Concurrency; i++ {
		go route.handle()
	}
}

func (route *Route) handle() {
	fromComp := route.components[route.From]

	for {
		select {
		case <-route.Ctx.Done():
			return
		case msg := <-fromComp.Inbound():
			for _, step := range route.steps {
				err := step(msg)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func (route *Route) Stop() {
	route.cancel()
}

func AsRoute(f any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			f,
			fx.ResultTags(`group:"nimble-routes"`),
		),
	)
}
