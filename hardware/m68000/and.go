package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type and struct {
	length  int
	size    instSize
	src     addressModeInfo
	destReg byte
	opmode  byte
}

func CreateAND(value uint16, code []uint8, codeOffset int) *and {
	opmode := valueFrom3Bits(value, 6)
	size := sizeFromOpmode(opmode)
	src, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset, codeOffset)
	destReg := valueFrom3Bits(value, 9)

	return &and{
		length:  2 + instSize,
		size:    size,
		src:     src,
		destReg: destReg,
		opmode:  opmode,
	}
}

func (a *and) Name() string {
	if a.opmode == 0 || a.opmode == 0b00000001 || a.opmode == 0b00000010 {
		return fmt.Sprintf(
			"and.%s   %s, D%d",
			a.size.string(),
			a.src.string(),
			a.destReg,
		)
	}

	return fmt.Sprintf(
		"and.%s   D%d, %s",
		a.size.string(),
		a.destReg,
		a.src.string())
}

func (a *and) Length() int {
	return a.length
}

func (a *and) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
