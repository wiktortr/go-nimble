package main

import (
	"github.com/wiktortr/go-nimble/components"
	"github.com/wiktortr/go-nimble/nimble"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func routeA() *nimble.Route {
	return nimble.
		From("timer://test?interval=1s").
		Process(func(message *nimble.Message) error {
			message.Payload = "test"
			return nil
		}).
		Log("Message from A").
		To("seda://test")
}

func routeB() *nimble.Route {
	return nimble.
		From("seda:test?buffSize=2").
		Log("Message from B")
}

func main() {
	fx.New(
		fx.Provide(zap.NewDevelopment),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		components.Core,
		nimble.AsRoute(routeA),
		nimble.AsRoute(routeB),
		nimble.Module,
	).Run()
}
