package timer

import (
	"context"
	"github.com/google/uuid"
	"github.com/wiktortr/go-nimble/nimble"
	"time"
)

type timerImpl struct {
	channel chan *nimble.Message
	dur     time.Duration
	cancel  context.CancelFunc
}

func (t *timerImpl) Start() {
	ticker := time.NewTicker(t.dur)
	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.channel <- &nimble.Message{
				Id: uuid.NewString(),
			}
		}
	}

}

func (t *timerImpl) Stop() {
	t.cancel()
}

func (t *timerImpl) Inbound() chan *nimble.Message {
	return t.channel
}

func (t *timerImpl) Process(*nimble.Message) error {
	return nil
}
