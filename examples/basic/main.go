package main

import (
	"github.com/wiktortr/go-nimble/components"
	"github.com/wiktortr/go-nimble/nimble"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func RouteA() *nimble.Route {
	return nimble.
		From("timer:test?interval=1s").
		Process(func(message *nimble.Message) error {
			message.Payload = "test"
			return nil
		}).
		Log("Message from A").
		To("seda:test")
}

func RouteB() *nimble.Route {
	return nimble.
		From("seda:test").
		Log("Message from B")
}

func main() {
	fx.New(
		fx.Provide(zap.NewDevelopment),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		components.Core,
		nimble.AsRoute(RouteA),
		nimble.AsRoute(RouteB),
		nimble.Module,
	).Run()
}
