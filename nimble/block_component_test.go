package nimble

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestComponentBlock_GetParent_ReturnsNil(t *testing.T) {
	block := &ComponentBlock{}
	assert.Nil(t, block.GetParent())
}

func TestComponentBlock_AddBlock_AlwaysReturnsError(t *testing.T) {
	// given
	block := &ComponentBlock{}

	// when
	err := block.AddBlock(&ComponentBlock{})

	// then
	assert.EqualError(t, err, "cannot add block to ComponentBlock")
}

func TestComponentBlock_Compile_ReturnsError_WhenRegistryInstantiateReturnsError(t *testing.T) {
	// given
	block := &ComponentBlock{Uri: "some-uri"}
	reg := &mockRegistry{}
	reg.On("Instantiate", mock.Anything).Return(&mockComponentImpl{}, errors.New("fail"))

	// when
	_, err := block.Compile(reg)

	// then
	assert.EqualError(t, err, "fail")
}

func TestComponentBlock_Compile_ReturnsProcessor_WhenRegistryInstantiateSucceeds(t *testing.T) {
	// given
	proc := &mockComponentImpl{}
	proc.On("Process", mock.AnythingOfType("*nimble.Message")).Return(nil)
	reg := &mockRegistry{}
	reg.On("Instantiate", mock.Anything).Return(proc, nil)
	block := &ComponentBlock{Uri: "some-uri"}

	// when
	processor, err := block.Compile(reg)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, processor)
	err = processor(&Message{})
	assert.NoError(t, err)
	proc.AssertCalled(t, "Process", mock.AnythingOfType("*nimble.Message"))
}

func TestComponentBlock_Compile_ProcessorReturnsError(t *testing.T) {
	// given
	proc := &mockComponentImpl{}
	proc.On("Process", mock.AnythingOfType("*nimble.Message")).Return(errors.New("fail"))
	reg := &mockRegistry{}
	reg.On("Instantiate", mock.Anything).Return(proc, nil)
	block := &ComponentBlock{Uri: "some-uri"}

	// when
	processor, err := block.Compile(reg)

	// then
	assert.NoError(t, err)
	err = processor(&Message{})
	assert.EqualError(t, err, "fail")
}
