package instructions

type nop struct {
}

func CreateNOP(value uint16, code []uint8, codeOffset int) *nop {
	return &nop{}
}

func (n *nop) Name() string {
	return "nop"
}

func (n *nop) Length() int {
	return 2
}
