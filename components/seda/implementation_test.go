package seda

import (
	"github.com/stretchr/testify/assert"
	"github.com/wiktortr/go-nimble/nimble"
	"testing"
)

func TestSedaImpl_Inbound_ReturnsChannel(t *testing.T) {
	// given
	ch := make(chan *nimble.Message, 1)
	impl := &sedaImpl{channel: ch}
	// when
	result := impl.Inbound()
	// then
	assert.Equal(t, ch, result)
}

func TestSedaImpl_Process_SendsMessageToChannel(t *testing.T) {
	// given
	ch := make(chan *nimble.Message, 1)
	impl := &sedaImpl{channel: ch}
	msg := nimble.NewMessage(nil)
	// when
	err := impl.Process(msg)
	// then
	assert.NoError(t, err)
	assert.Equal(t, msg, <-ch)
}
