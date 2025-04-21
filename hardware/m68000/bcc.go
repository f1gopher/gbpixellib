package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type bcc struct {
	size         instSize
	displacement int
	condition    testCondition

	// TODO - only needed for debug output, remove later and use PC value
	tmpStartOffset int
}

func CreateBcc(value uint16, code []uint8, codeOffset int) *bcc {
	size := zeroSize
	var displacement int = 0
	smallDisplacement := uint8(value & 0b0000000011111111)
	if smallDisplacement == 0 {
		size = wordSize
		displacement = twosComplementToInt16Bit(uint16(readData(code, codeOffset, size)))
	} else {
		displacement = twosComplementToInt8Bit(smallDisplacement)
	}

	return &bcc{
		size:           size,
		displacement:   displacement,
		condition:      condition(value),
		tmpStartOffset: codeOffset,
	}
}

func (b *bcc) Name() string {
	return fmt.Sprintf(
		"b%s     $%x",
		b.condition.string(),
		b.tmpStartOffset+b.displacement,
	)
}

func (b *bcc) Length() int {
	return 2 + b.size.bytesUsed()
}

func (b *bcc) Step(registers *Registers, memory memory.Memory[uint32]) bool {
	return true
}
