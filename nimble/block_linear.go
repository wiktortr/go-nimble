package nimble

type LinearBlock struct {
	Parent Block
	Blocks []Block
}

func (m *LinearBlock) GetParent() Block {
	return m.Parent
}

func (m *LinearBlock) AddBlock(block Block) error {
	m.Blocks = append(m.Blocks, block)
	return nil
}

func (m *LinearBlock) Compile(reg Registry) (MsgProcessor, error) {
	compiled := func(msg *Message) error {
		return nil
	}

	for i := len(m.Blocks) - 1; i >= 0; i-- {
		blockFunc, err := m.Blocks[i].Compile(reg)
		if err != nil {
			return nil, err
		}

		prev := compiled

		compiled = func(msg *Message) error {
			if err := blockFunc(msg); err != nil {
				return err
			}
			return prev(msg)
		}
	}

	return compiled, nil
}
