package nimble

type Block interface {
	GetParent() Block
	AddBlock(block Block) error
	Compile(reg *Registry) (MsgProcessor, error)
}
