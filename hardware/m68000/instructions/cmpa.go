package instructions

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type cmpa struct {
	length int
	size   instSize
	src    addressModeInfo
	reg    byte
}

func CreateCMPA(value uint16, code []uint8, codeOffset int) *cmpa {
	src, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)
	reg := valueFrom3Bits(value, 9)
	opmode := valueFrom3Bits(value, 6)
	size := wordSize
	if opmode == 0b00000111 {
		size = longSize
	}
	return &cmpa{
		length: 2 + instSize,
		size:   size,
		src:    src,
		reg:    reg,
	}
}

func (c *cmpa) Name() string {
	return fmt.Sprintf(
		"cmpa.%s  %s, A%d",
		c.size.string(),
		c.src.string(),
		c.reg)
}

func (c *cmpa) Length() int {
	return c.length
}

func (c *cmpa) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
