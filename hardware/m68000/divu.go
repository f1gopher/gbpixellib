package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type divu struct {
	length int
	size   instSize
	src    addressModeInfo
	reg    byte
}

func CreateDIVU(value uint16, code []uint8, codeOffset int) *divu {
	size := wordSize
	reg := valueFrom3Bits(value, 9)
	src, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset, codeOffset)
	return &divu{
		length: 2 + instSize,
		size:   size,
		src:    src,
		reg:    reg,
	}
}

func (d *divu) Name() string {
	return fmt.Sprintf(
		"divu.%s  %s, D%d",
		d.size.string(),
		d.src.string(),
		d.reg)
}

func (d *divu) Length() int {
	return d.length
}

func (d *divu) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
