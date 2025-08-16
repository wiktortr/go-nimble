package seda

import "github.com/wiktortr/go-nimble/nimble"

type sedaImpl struct {
	channel chan *nimble.Message
}

func (m *sedaImpl) Start() {

}

func (m *sedaImpl) Stop() {
}

func (m *sedaImpl) Inbound() chan *nimble.Message {
	return m.channel
}

func (m *sedaImpl) Process(msg *nimble.Message) error {
	m.channel <- msg
	return nil
}
