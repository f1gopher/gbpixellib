package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type not struct {
	length int
	size   instSize
	src    addressModeInfo
}

func CreateNOT(value uint16, code []uint8, codeOffset int) *not {
	size := twoBitSize(value, 6)
	src, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset, codeOffset)

	if src.mode == addressRegisterDirect ||
		src.mode == programCounterIndirectWithDisplacement ||
		src.mode == preIndexedImediate ||
		src.mode == programCounterIndirectWithIndex8bitDisplacement ||
		src.mode == invalidAddressMode {
		return nil
	}

	return &not{
		length: 2 + instSize,
		size:   size,
		src:    src,
	}
}

func (n *not) Name() string {
	return fmt.Sprintf(
		"not.%s   %s",
		n.size.string(),
		n.src.string())
}

func (n *not) Length() int {
	return n.length
}

func (n *not) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
