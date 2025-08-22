package nimble

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBlockFilter_GetParent_ReturnsCorrectParent(t *testing.T) {
	parent := &mockBlock{}
	block := &FilterBlock{Parent: parent}
	assert.Equal(t, parent, block.GetParent())
}

func TestBlockFilter_AddBlock_DelegatesToInnerBlockAndReturnsError(t *testing.T) {
	inner := &mockBlock{}
	block := &FilterBlock{InnerBlock: inner}
	expectedErr := errors.New("add error")
	inner.On("AddBlock", mock.Anything).Return(expectedErr)
	err := block.AddBlock(&mockBlock{})
	assert.Equal(t, expectedErr, err)
	inner.AssertExpectations(t)
}

func TestBlockFilter_AddBlock_DelegatesToInnerBlockAndReturnsNil(t *testing.T) {
	inner := &mockBlock{}
	block := &FilterBlock{InnerBlock: inner}
	inner.On("AddBlock", mock.Anything).Return(nil)
	err := block.AddBlock(&mockBlock{})
	assert.NoError(t, err)
	inner.AssertExpectations(t)
}

func TestBlockFilter_Compile_ReturnsErrorWhenInnerBlockCompileFails(t *testing.T) {
	inner := &mockBlock{}
	block := &FilterBlock{InnerBlock: inner}
	reg := &mockRegistry{}
	inner.On("Compile", reg).Return(nil, errors.New("compile error"))
	_, err := block.Compile(reg)
	assert.EqualError(t, err, "compile error")
	inner.AssertExpectations(t)
}

func TestBlockFilter_Compile_FiltersTrue_CallsInnerBlockProcessor(t *testing.T) {
	called := false
	reg := &mockRegistry{}
	inner := &mockBlock{}
	var blockFunc MsgProcessor = func(msg *Message) error { called = true; return nil }
	inner.On("Compile", reg).Return(blockFunc, nil)
	block := &FilterBlock{
		Filter:     func(msg *Message) (bool, error) { return true, nil },
		InnerBlock: inner,
	}

	proc, err := block.Compile(reg)
	assert.NoError(t, err)
	err = proc(&Message{})
	assert.NoError(t, err)
	assert.True(t, called)
	//inner.AssertExpectations(t)
}

func TestBlockFilter_Compile_FiltersFalse_DoesNotCallInnerBlockProcessor(t *testing.T) {
	reg := &mockRegistry{}
	inner := &mockBlock{}
	block := &FilterBlock{
		Filter:     func(msg *Message) (bool, error) { return false, nil },
		InnerBlock: inner,
	}
	var blockFunc MsgProcessor = func(msg *Message) error { return errors.New("should not be called") }
	inner.On("Compile", reg).Return(blockFunc, nil)

	proc, err := block.Compile(reg)
	assert.NoError(t, err)
	err = proc(&Message{})
	assert.NoError(t, err)
	inner.AssertExpectations(t)
}

func TestBlockFilter_Compile_FilterReturnsError_ReturnsError(t *testing.T) {
	inner := &mockBlock{}
	block := &FilterBlock{
		Filter:     func(msg *Message) (bool, error) { return false, errors.New("filter error") },
		InnerBlock: inner,
	}
	reg := &mockRegistry{}
	var blockFunc MsgProcessor = func(msg *Message) error { return errors.New("should not be called") }
	inner.On("Compile", reg).Return(blockFunc, nil)

	proc, err := block.Compile(reg)
	assert.NoError(t, err)
	err = proc(&Message{})
	assert.EqualError(t, err, "filter error")
	inner.AssertExpectations(t)
}
