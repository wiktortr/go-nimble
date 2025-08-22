package nimble

import "go.uber.org/zap"

type FilterFunc func(msg *Message) (bool, error)

func (m *Route) Filter(filter FilterFunc) *Route {
	block := &FilterBlock{
		Parent:     m.currentBlock,
		Filter:     filter,
		InnerBlock: &LinearBlock{m.currentBlock, []Block{}},
	}
	err := m.currentBlock.AddBlock(block)
	if err != nil {
		m.Registry.Logger().Error("Failed to add block", zap.Error(err), zap.String("route", m.Name))
	}
	m.currentBlock = block
	return m
}
