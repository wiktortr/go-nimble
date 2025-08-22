package nimble

type FilterBlock struct {
	Parent     Block
	Filter     FilterFunc
	InnerBlock Block
}

func (m *FilterBlock) GetParent() Block {
	return m.Parent
}

func (m *FilterBlock) AddBlock(block Block) error {
	return m.InnerBlock.AddBlock(block)
}

func (m *FilterBlock) Compile(reg Registry) (MsgProcessor, error) {
	comp, err := m.InnerBlock.Compile(reg)
	if err != nil {
		return nil, err
	}

	return func(msg *Message) error {
		val, err := m.Filter(msg)
		if err != nil {
			return err
		}
		if val {
			return comp(msg)
		}
		return nil
	}, nil
}
