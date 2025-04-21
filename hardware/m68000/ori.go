package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type ori struct {
	length int
	size   instSize
	dest   addressModeInfo

	data uint32
}

func CreateORI(value uint16, code []uint8, codeOffset int) *ori {
	size := twoBitSize(value, 6)
	dest, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset+size.bytesUsed(), codeOffset)

	if dest.mode == addressRegisterDirect ||
		dest.mode == programCounterIndirectWithDisplacement ||
		dest.mode == preIndexedImediate ||
		dest.mode == programCounterIndirectWithIndex8bitDisplacement ||
		dest.mode == invalidAddressMode ||
		size == zeroSize {
		return nil
	}

	return &ori{
		length: 2 + instSize + size.bytesUsed(),
		size:   size,
		dest:   dest,
		data:   readData(code, codeOffset, size),
	}
}

func (o *ori) Name() string {
	return fmt.Sprintf(
		"ori.%s   #$%x, %s",
		o.size.string(),
		o.data,
		o.dest.string())
}

func (o *ori) Length() int {
	return o.length
}

func (o *ori) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
