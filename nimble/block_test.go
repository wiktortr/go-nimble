package nimble

import "github.com/stretchr/testify/mock"

type mockBlock struct {
	mock.Mock
}

func (m *mockBlock) AddBlock(block Block) error {
	args := m.Called(block)
	return args.Error(0)
}
func (m *mockBlock) Compile(reg Registry) (MsgProcessor, error) {
	args := m.Called(reg)
	if args.Get(0) != nil {
		return args.Get(0).(MsgProcessor), nil
	}
	return nil, args.Error(1)
}

func (m *mockBlock) GetParent() Block {
	args := m.Called()
	return args.Get(0).(Block)
}
