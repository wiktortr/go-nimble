package timer

import (
	"github.com/stretchr/testify/assert"
	"github.com/wiktortr/go-nimble/nimble"
	"testing"
	"time"
)

func TestTimer_Key_ReturnsCorrectKey(t *testing.T) {
	// given
	timer := &Timer{}
	// when
	key := timer.Key()
	// then
	assert.Equal(t, "timer", key)
}

func TestTimer_Instantiate_ReturnsImplWithValidInterval(t *testing.T) {
	// given
	timer := &Timer{}
	params := &nimble.ComponentParams{
		Values: map[string][]string{"interval": {"2s"}},
	}
	// when
	inst, err := timer.Instantiate(nil, params)
	// then
	assert.NoError(t, err)
	impl, ok := inst.(*timerImpl)
	assert.True(t, ok)
	assert.Equal(t, 2*time.Second, impl.dur)
	assert.NotNil(t, impl.channel)
}

func TestTimer_Instantiate_InvalidInterval_ReturnsError(t *testing.T) {
	// given
	timer := &Timer{}
	params := &nimble.ComponentParams{
		Values: map[string][]string{"interval": {"invalid"}},
	}
	// when
	inst, err := timer.Instantiate(nil, params)
	// then
	assert.Nil(t, inst)
	assert.Error(t, err)
}

func TestTimer_Instantiate_MissingInterval_ReturnsError(t *testing.T) {
	// given
	timer := &Timer{}
	params := &nimble.ComponentParams{
		Values: map[string][]string{},
	}
	// when
	inst, err := timer.Instantiate(nil, params)
	// then
	assert.Nil(t, inst)
	assert.Error(t, err)
}
