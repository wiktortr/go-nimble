package autostart

import (
	"github.com/stretchr/testify/assert"
	"github.com/wiktortr/go-nimble/nimble"
	"testing"
	"time"
)

func TestAutostartImpl_Start_SendsMessageAfterDelay(t *testing.T) {
	// given
	impl := &autostartImpl{
		delay:   10 * time.Millisecond,
		channel: make(chan *nimble.Message, 1),
	}
	// when
	impl.Start()
	// then
	select {
	case msg := <-impl.channel:
		assert.NotNil(t, msg)
	case <-time.After(50 * time.Millisecond):
		t.Error("Message was not sent after delay")
	}
}

func TestAutostartImpl_Stop_DoesNothing(t *testing.T) {
	// given
	impl := &autostartImpl{
		delay:   1 * time.Millisecond,
		channel: make(chan *nimble.Message, 1),
	}
	// when
	impl.Stop()
	// then
	// brak efektów ubocznych, test przechodzi jeśli nie ma paniki
}

func TestAutostartImpl_Inbound_ReturnsChannel(t *testing.T) {
	// given
	ch := make(chan *nimble.Message, 1)
	impl := &autostartImpl{
		delay:   1 * time.Millisecond,
		channel: ch,
	}
	// when
	result := impl.Inbound()
	// then
	assert.Equal(t, ch, result)
}

func TestAutostartImpl_Process_AlwaysReturnsNil(t *testing.T) {
	// given
	impl := &autostartImpl{
		delay:   1 * time.Millisecond,
		channel: make(chan *nimble.Message, 1),
	}
	msg := nimble.NewMessage(nil)
	// when
	err := impl.Process(msg)
	// then
	assert.NoError(t, err)
}
