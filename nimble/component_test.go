package nimble

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"
	"testing"
)

type mockComponentImpl struct {
	mock.Mock
}

func (m *mockComponentImpl) Start() {}

func (m *mockComponentImpl) Stop() {}

func (m *mockComponentImpl) Inbound() chan *Message {
	return m.Called().Get(0).(chan *Message)
}

func (m *mockComponentImpl) Process(msg *Message) error {
	args := m.Called(msg)
	return args.Error(0)
}

type mockComponent struct {
	mock.Mock
}

func (m *mockComponent) Key() string {
	return m.Called().String(0)
}

func (m *mockComponent) Instantiate(reg Registry, params *ComponentParams) (ComponentImpl, error) {
	args := m.Called(reg, params)
	return args.Get(0).(ComponentImpl), args.Error(1)
}

func TestComponent_Key_ReturnsExpectedKey(t *testing.T) {
	m := &mockComponent{}
	m.On("Key").Return("test-key")
	assert.Equal(t, "test-key", m.Key())
	m.AssertExpectations(t)
}

func TestComponent_Instantiate_ReturnsComponentImplAndError(t *testing.T) {
	m := &mockComponent{}
	params := &ComponentParams{}
	reg := &mockRegistry{}
	impl := &mockComponentImpl{}
	m.On("Instantiate", reg, params).Return(impl, nil)
	result, err := m.Instantiate(reg, params)
	assert.NoError(t, err)
	assert.Equal(t, impl, result)
	m.AssertExpectations(t)
}

func TestComponent_Instantiate_ReturnsError(t *testing.T) {
	m := &mockComponent{}
	params := &ComponentParams{}
	reg := &mockRegistry{}
	m.On("Instantiate", reg, params).Return(&mockComponentImpl{}, assert.AnError)
	_, err := m.Instantiate(reg, params)
	assert.Error(t, err)
	m.AssertExpectations(t)
}

func TestComponent_AsComponent_ReturnsFxOption(t *testing.T) {
	opt := AsComponent(func() *mockComponent { return &mockComponent{} })
	_, ok := opt.(fx.Option)
	assert.True(t, ok)
}
