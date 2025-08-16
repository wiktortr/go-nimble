package main

import (
	"github.com/wiktortr/go-nimble/components"
	"github.com/wiktortr/go-nimble/nimble"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"math/rand"
)

func graterThan5(msg *nimble.Message) (bool, error) {
	return msg.Payload.(int) > 5, nil
}

func RouteA() *nimble.Route {
	r := rand.New(rand.NewSource(99))
	return nimble.
		From("timer:test?interval=1s").
		Process(func(message *nimble.Message) error {
			message.Payload = r.Intn(10)
			return nil
		}).
		Log("Message payload").
		Filter(graterThan5).
		Log("Message filter").
		End().
		Log("After filter block")
}

func main() {
	fx.New(
		fx.Provide(zap.NewDevelopment),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		components.Core,
		nimble.AsRoute(RouteA),
		nimble.Module,
	).Run()
}
