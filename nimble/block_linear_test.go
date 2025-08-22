package nimble

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

var happyMsgProc MsgProcessor = func(msg *Message) error { return nil }
var failMsgProc MsgProcessor = func(msg *Message) error { return errors.New("fail") }

func TestLinearBlock_GetParent_ReturnsParent(t *testing.T) {
	parent := &mockBlock{}
	block := &LinearBlock{Parent: parent}
	assert.Equal(t, parent, block.GetParent())
}

func TestLinearBlock_AddBlock_AppendsBlock(t *testing.T) {
	block := &LinearBlock{}
	child := &mockBlock{}
	err := block.AddBlock(child)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(block.Blocks))
	assert.Equal(t, child, block.Blocks[0])
}

func TestLinearBlock_Compile_NoBlocks_ReturnsNoopProcessor(t *testing.T) {
	// given
	block := &LinearBlock{}

	// when
	proc, err := block.Compile(nil)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, proc)
	assert.NoError(t, proc(&Message{}))
}

func TestLinearBlock_Compile_AllBlocksReturnNil_NoError(t *testing.T) {
	// given
	block1 := &mockBlock{}
	block2 := &mockBlock{}
	reg := &mockRegistry{}
	block1.On("Compile", reg).Return(happyMsgProc, nil)
	block2.On("Compile", reg).Return(happyMsgProc, nil)
	linear := &LinearBlock{Blocks: []Block{block1, block2}}

	// when
	proc, err := linear.Compile(reg)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, proc)
	assert.NoError(t, proc(&Message{}))
	block1.AssertCalled(t, "Compile", reg)
	block2.AssertCalled(t, "Compile", reg)
}

func TestLinearBlock_Compile_BlockReturnsError_PropagatesError(t *testing.T) {
	// given
	expectedErr := errors.New("compile error")
	reg := &mockRegistry{}
	block1 := &mockBlock{}
	block1.On("Compile", reg).Return(nil, expectedErr)
	block2 := &mockBlock{}
	block2.On("Compile", reg).Return(happyMsgProc, nil)
	linear := &LinearBlock{Blocks: []Block{block1, block2}}

	// when
	proc, err := linear.Compile(reg)

	// then
	assert.Nil(t, proc)
	assert.Equal(t, expectedErr, err)
	block1.AssertExpectations(t)
}

func TestLinearBlock_Compile_ProcessorReturnsError_StopsChain(t *testing.T) {
	// given
	reg := &mockRegistry{}
	block1 := &mockBlock{}
	block1.On("Compile", reg).Return(failMsgProc, nil)
	block2 := &mockBlock{}
	block2.On("Compile", reg).Return(happyMsgProc, nil)
	linear := &LinearBlock{Blocks: []Block{block2, block1}}

	// when
	proc, err := linear.Compile(reg)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, proc)
	err = proc(&Message{})
	assert.EqualError(t, err, "fail")
	block1.AssertExpectations(t)
	block2.AssertExpectations(t)
}
