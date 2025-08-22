package seda

import (
	"github.com/stretchr/testify/assert"
	"github.com/wiktortr/go-nimble/nimble"
	"testing"
)

func TestSeda_Key_ReturnsCorrectKey(t *testing.T) {
	// given
	seda := &Seda{}
	// when
	key := seda.Key()
	// then
	assert.Equal(t, "seda", key)
}

func TestSeda_Instantiate_ReturnsImplWithDefaultBuffSize(t *testing.T) {
	// given
	seda := &Seda{}
	params := &nimble.ComponentParams{
		Values: map[string][]string{},
	}
	// when
	inst, err := seda.Instantiate(nil, params)
	// then
	assert.NoError(t, err)
	impl, ok := inst.(*sedaImpl)
	assert.True(t, ok)
	assert.NotNil(t, impl.channel)
	assert.Equal(t, 100, cap(impl.channel))
}

func TestSeda_Instantiate_ReturnsImplWithCustomBuffSize(t *testing.T) {
	// given
	seda := &Seda{}
	params := &nimble.ComponentParams{
		Values: map[string][]string{"buffSize": {"42"}},
	}
	// when
	inst, err := seda.Instantiate(nil, params)
	// then
	assert.NoError(t, err)
	impl, ok := inst.(*sedaImpl)
	assert.True(t, ok)
	assert.NotNil(t, impl.channel)
	assert.Equal(t, 42, cap(impl.channel))
}

func TestSeda_Instantiate_InvalidBuffSize_UsesDefault(t *testing.T) {
	// given
	seda := &Seda{}
	params := &nimble.ComponentParams{
		Values: map[string][]string{"buffSize": {"invalid"}},
	}
	// when
	inst, err := seda.Instantiate(nil, params)
	// then
	assert.NoError(t, err)
	impl, ok := inst.(*sedaImpl)
	assert.True(t, ok)
	assert.NotNil(t, impl.channel)
	assert.Equal(t, 100, cap(impl.channel))
}
