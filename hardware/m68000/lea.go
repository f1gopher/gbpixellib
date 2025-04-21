package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type lea struct {
	length int
	size   instSize
	src    addressModeInfo
	dstReg byte
	data   uint32
}

func CreateLEA(value uint16, code []uint8, codeOffset int) *lea {
	src, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)
	size := wordSize
	if src.mode == absoluteLong {
		size = longSize
	}

	dstReg := byte((value >> 9) & 0b0000000000000111)

	if src.mode == invalidAddressMode {
		return nil
	}

	return &lea{
		length: 2 + instSize,
		size:   size,
		src:    src,
		data:   readData(code, codeOffset, size),
		dstReg: dstReg,
	}
}

func (l *lea) Name() string {
	return fmt.Sprintf("lea     %s, A%d",
		l.src.string(),
		l.dstReg)
}

func (l *lea) Length() int {
	return l.length
}

func (l *lea) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
