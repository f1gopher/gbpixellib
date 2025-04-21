package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type bra struct {
	size         instSize
	displacement int

	tmpStartOffset int
}

func CreateBRA(value uint16, code []uint8, codeOffset int) *bra {
	size := zeroSize
	var displacement int = 0
	smallDisplacement := uint8(value & 0b0000000011111111)
	if smallDisplacement == 0 {
		size = wordSize
		displacement = twosComplementToInt16Bit(uint16(readData(code, codeOffset, size)))
	} else {
		displacement = twosComplementToInt8Bit(smallDisplacement)
	}

	return &bra{
		size:           size,
		displacement:   displacement,
		tmpStartOffset: codeOffset,
	}
}

func (b *bra) Name() string {
	return fmt.Sprintf("bra     $%x", b.tmpStartOffset+b.displacement)
}

func (b *bra) Length() int {
	return 2 + b.size.bytesUsed()
}

func (b *bra) Step(registers *Registers, memory memory.Memory[uint32]) bool {
	return true
}
