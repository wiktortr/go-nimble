package timer

import (
	"github.com/stretchr/testify/assert"
	"github.com/wiktortr/go-nimble/nimble"
	"testing"
	"time"
)

func TestTimerImpl_Start_SendsMessagesAtIntervalAndStopsOnCancel(t *testing.T) {
	// given
	ch := make(chan *nimble.Message, 2)
	impl := &timerImpl{
		channel: ch,
		dur:     10 * time.Millisecond,
	}
	done := make(chan struct{})
	go func() {
		impl.Start()
		close(done)
	}()
	// when
	time.Sleep(25 * time.Millisecond)
	impl.Stop()
	<-done
	// then
	count := 0
loop:
	for {
		select {
		case <-ch:
			count++
		default:
			break loop
		}
	}
	assert.GreaterOrEqual(t, count, 2)
}

func TestTimerImpl_Stop_CancelsContextWithoutPanic(t *testing.T) {
	// given
	impl := &timerImpl{
		channel: make(chan *nimble.Message),
		dur:     10 * time.Millisecond,
	}
	impl.cancel = func() {}
	// when
	impl.Stop()
	// then
	// brak paniki, test przechodzi
}

func TestTimerImpl_Inbound_ReturnsChannel(t *testing.T) {
	// given
	ch := make(chan *nimble.Message)
	impl := &timerImpl{channel: ch}
	// when
	result := impl.Inbound()
	// then
	assert.Equal(t, ch, result)
}

func TestTimerImpl_Process_AlwaysReturnsNil(t *testing.T) {
	// given
	impl := &timerImpl{channel: make(chan *nimble.Message)}
	msg := nimble.NewMessage(nil)
	// when
	err := impl.Process(msg)
	// then
	assert.NoError(t, err)
}
