package instructions

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type bsr struct {
	length       int
	displacement int

	tmpStartOffset int
}

func CreateBSR(value uint16, code []uint8, codeOffset int) *bsr {
	length := 2
	displacement := twosComplementToInt8Bit(uint8(value & 0b0000000011111111))
	if value&0b0000000011111111 == 0 {
		length += 2
		displacement = twosComplementToInt16Bit(uint16(readData(code, codeOffset, wordSize)))
	}

	return &bsr{
		length:         length,
		displacement:   displacement,
		tmpStartOffset: codeOffset,
	}
}

func (b *bsr) Name() string {
	return fmt.Sprintf(
		"bsr     $%x",
		b.tmpStartOffset+b.displacement)
}

func (b *bsr) Length() int {
	return b.length
}

func (b *bsr) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
