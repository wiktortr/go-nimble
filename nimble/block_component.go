package nimble

import "errors"

type ComponentBlock struct {
	Uri string
}

func (m *ComponentBlock) GetParent() Block {
	return nil
}

func (m *ComponentBlock) AddBlock(_ Block) error {
	return errors.New("cannot add block to ComponentBlock")
}

func (m *ComponentBlock) Compile(reg *Registry) (MsgProcessor, error) {
	instantiate, err := reg.Instantiate(m.Uri)
	if err != nil {
		return nil, err
	}
	return func(msg *Message) error {
		return instantiate.Process(msg)
	}, nil
}
