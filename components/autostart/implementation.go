package autostart

import (
	"github.com/wiktortr/go-nimble/nimble"
	"time"
)

type autostartImpl struct {
	delay   time.Duration
	channel chan *nimble.Message
}

func (m *autostartImpl) Start() {
	timer := time.NewTimer(m.delay)
	go func() {
		<-timer.C
		m.channel <- nimble.NewMessage(nil)
	}()
}

func (m *autostartImpl) Stop() {}

func (m *autostartImpl) Inbound() chan *nimble.Message {
	return m.channel
}

func (m *autostartImpl) Process(_ *nimble.Message) error {
	return nil
}
