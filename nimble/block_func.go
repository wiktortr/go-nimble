package nimble

import "errors"

type FunctionalBlock struct {
	Function MsgProcessor
}

func (m *FunctionalBlock) GetParent() Block {
	return nil
}

func (m *FunctionalBlock) AddBlock(block Block) error {
	return errors.New("cannot add block to FunctionalBlock")
}

func (m *FunctionalBlock) Compile(_ *Registry) (MsgProcessor, error) {
	return m.Function, nil
}
