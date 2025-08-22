package nimble

import "go.uber.org/zap"

func (m *Route) Process(process MsgProcessor) *Route {
	err := m.currentBlock.AddBlock(&FunctionalBlock{process})
	if err != nil {
		m.Registry.Logger().Error("Failed to add block", zap.Error(err), zap.String("route", m.Name))
	}
	return m
}
