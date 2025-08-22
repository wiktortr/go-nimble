package autostart

import (
	"github.com/stretchr/testify/assert"
	"github.com/wiktortr/go-nimble/nimble"
	"testing"
	"time"
)

func TestAutoStart_Key_ReturnsCorrectKey(t *testing.T) {
	// given
	auto := &AutoStart{}
	// when
	key := auto.Key()
	// then
	assert.Equal(t, "autostart", key)
}

func TestAutoStart_Instantiate_ReturnsImplWithDefaultDelay(t *testing.T) {
	// given
	auto := &AutoStart{}
	params := &nimble.ComponentParams{
		Values: map[string][]string{},
	}
	// when
	inst, err := auto.Instantiate(nil, params)
	// then
	assert.NoError(t, err)
	impl, ok := inst.(*autostartImpl)
	assert.True(t, ok)
	assert.Equal(t, time.Second, impl.delay)
	assert.NotNil(t, impl.channel)
}

func TestAutoStart_Instantiate_ReturnsImplWithCustomDelay(t *testing.T) {
	// given
	auto := &AutoStart{}
	params := &nimble.ComponentParams{
		Values: map[string][]string{"delay": {"2s"}},
	}
	// when
	inst, err := auto.Instantiate(nil, params)
	// then
	assert.NoError(t, err)
	impl, ok := inst.(*autostartImpl)
	assert.True(t, ok)
	assert.Equal(t, 2*time.Second, impl.delay)
	assert.NotNil(t, impl.channel)
}

func TestAutoStart_Instantiate_ReturnsImplWithInvalidDelay_UsesDefault(t *testing.T) {
	// given
	auto := &AutoStart{}
	params := &nimble.ComponentParams{
		Values: map[string][]string{"delay": {"invalid"}},
	}
	// when
	inst, err := auto.Instantiate(nil, params)
	// then
	assert.NoError(t, err)
	impl, ok := inst.(*autostartImpl)
	assert.True(t, ok)
	assert.Equal(t, time.Second, impl.delay)
	assert.NotNil(t, impl.channel)
}
