package nimble

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFuncBlock_AddBlock_ReturnsError(t *testing.T) {
	block := &FunctionalBlock{Function: nil}
	err := block.AddBlock(&FunctionalBlock{})
	assert.EqualError(t, err, "cannot add block to FunctionalBlock")
}

func TestFuncBlock_GetParent_AlwaysReturnsNil(t *testing.T) {
	block := &FunctionalBlock{Function: nil}
	assert.Nil(t, block.GetParent())
}

func TestFuncBlock_Compile_ReturnsNilFunction_WhenFunctionIsNil(t *testing.T) {
	block := &FunctionalBlock{Function: nil}
	processor, err := block.Compile(nil)
	assert.NoError(t, err)
	assert.Nil(t, processor)
}

func TestFuncBlock_Compile_ReturnsFunction_WhenFunctionIsSet(t *testing.T) {
	called := false
	fn := func(msg *Message) error {
		called = true
		return nil
	}
	block := &FunctionalBlock{Function: fn}
	processor, err := block.Compile(nil)
	assert.NoError(t, err)
	assert.NotNil(t, processor)
	_ = processor(&Message{})
	assert.True(t, called)
}

func TestCompile_FunctionReturnsError_PropagatesError(t *testing.T) {
	expectedErr := errors.New("some error")
	fn := func(msg *Message) error {
		return expectedErr
	}
	block := &FunctionalBlock{Function: fn}
	processor, err := block.Compile(nil)
	assert.NoError(t, err)
	err = processor(&Message{})
	assert.Equal(t, expectedErr, err)
}
