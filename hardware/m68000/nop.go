package m68000

import "github.com/f1gopher/gbpixellib/hardware/memory"

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

func (n *nop) Step(registers *Registers, memory memory.Memory[uint32]) bool {
	return true
}
